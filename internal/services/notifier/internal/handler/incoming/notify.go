package incoming

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/message"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/broker"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	message_dto "github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_notifier_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
)

func HandleIncomingNotifyHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pBroker broker.IDataBroker,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeOFF)

		logBuilder := http_logger.NewLogBuilder(hls_notifier_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "method", api.WithoutContentType())
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

		aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName)
		if aliasName == "" {
			pLogger.PushErro(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadGateway, "alias", api.WithoutContentType())
			return
		}

		msg := message_dto.NewMessageBody(msgBody, time.Now())
		msg.SetIsNotify(true)

		pBroker.Produce(message.NewMessageContainer(aliasName, msg))

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, "ok", api.WithoutContentType())
	}
}
