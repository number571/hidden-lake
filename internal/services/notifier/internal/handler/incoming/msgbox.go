package incoming

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/message"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/broker"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	message_dto "github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/notifier/pkg/app/config"
	hls_notifier_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
)

func HandleIncomingMsgboxHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pDB database.IKVDatabase,
	pBroker broker.IDataBroker,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hlk_settings.CHeaderResponseMode, getResponseMode(pConfig))

		logBuilder := http_logger.NewLogBuilder(hls_notifier_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "method", api.WithoutContentType())
			return
		}

		aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName)
		if aliasName == "" {
			pLogger.PushErro(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadGateway, "alias", api.WithoutContentType())
			return
		}

		queryParams := pR.URL.Query()

		msgUID := queryParams.Get("uid")
		if !message_dto.MessageUIDIsValid(msgUID) {
			pLogger.PushWarn(logBuilder.WithMessage("msg_uid"))
			_ = api.Response(pW, http.StatusBadRequest, "uid", api.WithoutContentType())
			return
		}

		if pR.Method == http.MethodGet {
			msg, err := pDB.LoadByUID(aliasName, msgUID)
			if err == nil {
				pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
				_ = api.Response(pW, http.StatusOK, msg.ToBytes(), api.WithoutContentType())
				return
			}
			if errors.Is(err, database.ErrNotFoundMessage) {
				pLogger.PushWarn(logBuilder.WithMessage("message_not_found"))
				_ = api.Response(pW, http.StatusNotFound, "notfound", api.WithoutContentType())
				return
			}
			pLogger.PushErro(logBuilder.WithMessage("load_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "getmsg", api.WithoutContentType())
			return
		}

		msgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "readmsg", api.WithoutContentType())
			return
		}

		msgBody := string(msgBytes)
		if len(msgBody) == 0 || chars.HasNotGraphicCharacters(msgBody) {
			pLogger.PushWarn(logBuilder.WithMessage("has_not_graphic_chars"))
			_ = api.Response(pW, http.StatusBadRequest, "graph", api.WithoutContentType())
			return
		}

		_, err = pDB.LoadByUID(aliasName, msgUID)
		if err == nil {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusAccepted, "exist", api.WithoutContentType())
			return
		}
		if !errors.Is(err, database.ErrNotFoundMessage) {
			pLogger.PushErro(logBuilder.WithMessage("load_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "getmsg", api.WithoutContentType())
			return
		}

		msg := message_dto.NewMessage(
			message_dto.NewMessageHead(message_dto.CIncomingState, msgUID),
			message_dto.NewMessageBody(msgBody, time.Now()),
		)

		if _, err := pDB.Push(aliasName, msg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "setmsg", api.WithoutContentType())
			return
		}

		msgContainer := message.NewMessageContainer(
			aliasName,
			msg.GetBody(),
		).WithMessageUID(msgUID)

		pBroker.Produce(msgContainer)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, "ok", api.WithoutContentType())
	}
}

func getResponseMode(pConfig config.IConfig) string {
	if pConfig.GetSettings().GetResponseEnabled() {
		return hlk_settings.CHeaderResponseModeON
	}
	return hlk_settings.CHeaderResponseModeOFF
}
