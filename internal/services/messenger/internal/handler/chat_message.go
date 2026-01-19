package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/client/message"
	hls_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/pubkey"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func HandleChatMessageAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pHlkClient hlk_client.IClient,
	pDatabase database.IKVDatabase,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			size, err := getLimitOnPushRequestSize(pCtx, pHlkClient)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, strconv.FormatUint(size, 10))
			return
		}

		body, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		strBody := string(body)
		if len(strBody) == 0 || chars.HasNotGraphicCharacters(strBody) {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		aliasName := pR.URL.Query().Get("friend")
		fPubKey, err := pubkey.GetFriendPubKeyByAliasName(pCtx, pHlkClient, aliasName)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("load_pubkey"))
			_ = api.Response(pW, http.StatusForbidden, "failed: load public key")
			return
		}

		myPubKey, err := pHlkClient.GetPubKey(pCtx)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		if err := pHlkClient.SendRequest(pCtx, aliasName, newPushRequest(strBody)); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("send_request"))
			_ = api.Response(pW, http.StatusForbidden, "failed: send request")
			return
		}

		msg := message.NewMessage(false, string(body))
		if err := pDatabase.Push(database.NewRelation(myPubKey, fPubKey), msg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, msg.GetTimestamp())
	}
}

func newPushRequest(body string) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hls_settings.CAppShortName).
		WithPath(hls_settings.CPushPath).
		WithBody([]byte(body)).
		Build()
}

func getLimitOnPushRequestSize(pCtx context.Context, pHlkClient hlk_client.IClient) (uint64, error) {
	sett, err := pHlkClient.GetSettings(pCtx)
	if err != nil {
		return 0, err
	}

	reqSize := uint64(len(newPushRequest("").ToBytes()))
	pldLimit := sett.GetPayloadSizeBytes()
	if reqSize >= pldLimit {
		return 0, errors.New("request size >= payload limit") // nolint: err113
	}

	return pldLimit - reqSize, nil
}
