package handler

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/applications/messenger/internal/msgbroker"
	hlm_utils "github.com/number571/hidden-lake/internal/applications/messenger/internal/utils"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"

	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func HandleIncomingPushHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pDB database.IKVDatabase,
	pBroker msgbroker.IMessageBroker,
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeOFF)

		logBuilder := http_logger.NewLogBuilder(hlm_settings.CServiceName, pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		rawMsgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		fPubKey := asymmetric.LoadPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fPubKey == nil {
			pLogger.PushErro(logBuilder.WithMessage("load_pubkey"))
			_ = api.Response(pW, http.StatusForbidden, "failed: load public key")
			return
		}

		if err := isValidMsgBytes(rawMsgBytes); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("recv_message"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: get message bytes")
			return
		}

		myPubKey, err := pHlsClient.GetPubKey(pCtx)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		rel := database.NewRelation(myPubKey, fPubKey)
		dbMsg := database.NewMessage(true, rawMsgBytes)

		if err := pDB.Push(rel, dbMsg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		pBroker.Produce(
			fPubKey.GetHasher().ToString(),
			getMessage(
				true,
				dbMsg.GetMessage(),
				dbMsg.GetTimestamp(),
			),
		)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, hlm_settings.CServiceFullName)
	}
}

func isValidMsgBytes(rawMsgBytes []byte) error {
	switch {
	case isText(rawMsgBytes):
		strMsg := strings.TrimSpace(unwrapText(rawMsgBytes, true))
		if strMsg == "" {
			return ErrMessageNull
		}
		if chars.HasNotGraphicCharacters(strMsg) {
			return ErrHasNotWritableChars
		}
		return nil
	case isFile(rawMsgBytes):
		filename, msgBytes := unwrapFile(rawMsgBytes, true)
		if filename == "" || len(msgBytes) == 0 {
			return ErrUnwrapFile
		}
		return nil
	default:
		return ErrUnknownMessageType
	}
}

func getMessage(pEscape bool, pRawMsgBytes []byte, pTimestamp string) hlm_utils.SMessage {
	switch {
	case isText(pRawMsgBytes):
		msgData := unwrapText(pRawMsgBytes, pEscape)
		if msgData == "" {
			panic("message data = nil")
		}
		return hlm_utils.SMessage{
			FTimestamp: pTimestamp,
			FMainData:  msgData,
		}
	case isFile(pRawMsgBytes):
		filename, msgData := unwrapFile(pRawMsgBytes, pEscape)
		if filename == "" || msgData == "" {
			panic("filename = nil OR message data = nil")
		}
		return hlm_utils.SMessage{
			FFileName:  filename,
			FTimestamp: pTimestamp,
			FMainData:  msgData,
		}
	default:
		panic("got unknown message type")
	}
}
