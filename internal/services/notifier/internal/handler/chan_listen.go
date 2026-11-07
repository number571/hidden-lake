package handler

import (
	"context"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/message"
	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/broker"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleChanListenAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pBroker broker.IDataBroker,
) http.HandlerFunc {
	buildSettings := build.GetSettings()

	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		queryParams := pR.URL.Query()

		sid := queryParams.Get("sid")
		if err := pBroker.Register(sid); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("limit_subs"))
			_ = api.Response(pW, http.StatusNotAcceptable, "has limit of subscribers")
			return
		}

		ctx, cancel := context.WithTimeout(pCtx, buildSettings.GetHttpReadTimeout())
		defer cancel()

		v, err := pBroker.Consume(ctx, sid)
		if err != nil {
			pLogger.PushInfo(logBuilder.WithMessage("no_content"))
			_ = api.Response(pW, http.StatusNoContent, "")
			return
		}

		mc, ok := v.(message.IMessageContainer)
		if !ok {
			pLogger.PushErro(logBuilder.WithMessage("invalid_type"))
			_ = api.Response(pW, http.StatusInternalServerError, "invalid type")
			return
		}

		pW.Header().Set(hls_settings.CHeaderFriendName, mc.GetFriend())

		msgBody := mc.GetMessageBody()
		contentType := getContentType(msgBody)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, msgBody, api.WithCustomContentType(contentType))
	}
}

func getContentType(msgBody []byte) string {
	if chars.HasNotGraphicCharacters(string(msgBody)) {
		return api.CApplicationOctetStream
	}
	return api.CTextPlain
}
