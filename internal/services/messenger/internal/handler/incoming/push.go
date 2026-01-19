package incoming

import (
	"context"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/client/message"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/pubkey"

	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
)

func HandleIncomingPushHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pDB database.IKVDatabase,
	pBroker message.IMessageBroker,
	pHlkClient hlk_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeOFF)

		logBuilder := http_logger.NewLogBuilder(hls_messenger_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		msgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		rawMsg := string(msgBytes)
		if len(rawMsg) == 0 || chars.HasNotGraphicCharacters(rawMsg) {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: response message")
			return
		}

		aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName)
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

		msg := message.NewMessage(true, rawMsg)
		pBroker.Produce(aliasName, msg)

		if err := pDB.Push(database.NewRelation(myPubKey, fPubKey), msg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, http_logger.CLogSuccess)
	}
}
