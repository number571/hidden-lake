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
	hls_notifier_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
)

func HandleIncomingPushHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pDB database.IKVDatabase,
	pBroker broker.IDataBroker,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hls_notifier_settings.GetAppShortNameFMT(), pR)

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

		msgBody := string(msgBytes)
		if len(msgBody) == 0 || chars.HasNotGraphicCharacters(msgBody) {
			pLogger.PushWarn(logBuilder.WithMessage("has_not_graphic_chars"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: has not graphic characters")
			return
		}

		aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName)
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: read alias name")
			return
		}

		queryParams := pR.URL.Query()
		msgUID := queryParams.Get("uid")

		if !message_dto.MessageUIDIsValid(msgUID) {
			pLogger.PushErro(logBuilder.WithMessage("msg_uid"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: message uid")
			return
		}

		_, err = pDB.LoadByUID(aliasName, msgUID)
		if err == nil {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusAccepted, http_logger.CLogSuccess)
			return
		}
		if !errors.Is(err, database.ErrNotFoundMessage) {
			pLogger.PushErro(logBuilder.WithMessage("load_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: load message from database")
			return
		}

		msg, err := message_dto.NewMessageWithParams(
			message_dto.CIncomingState,
			msgBody,
			msgUID,
			time.Now(),
		)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("create_message"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: create message")
			return
		}
		if _, err := pDB.Push(aliasName, msg); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("push_message"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
			return
		}

		pBroker.Produce(message.NewMessageContainer(aliasName, msg))

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, http_logger.CLogSuccess)
	}
}
