package handler

import (
	"context"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/message"
	hls_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleChatSubscribeAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pBroker message.IMessageBroker,
) http.HandlerFunc {
	buildSettings := build.GetSettings()

	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		friend := pR.URL.Query().Get("friend")
		sid := pR.URL.Query().Get("sid")

		ctx, cancel := context.WithTimeout(pCtx, buildSettings.GetHttpReadTimeout())
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				pLogger.PushInfo(logBuilder.WithMessage("no_content"))
				_ = api.Response(pW, http.StatusNoContent, []byte{})
				return
			case c, ok := <-pBroker.Consume(sid):
				if !ok {
					pLogger.PushInfo(logBuilder.WithMessage("no_content"))
					_ = api.Response(pW, http.StatusNoContent, []byte{})
					return
				}
				if c.GetFriend() != friend {
					pLogger.PushInfo(logBuilder.WithMessage("another_friend"))
					continue
				}
				pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
				_ = api.Response(pW, http.StatusOK, c.GetMessage().ToString())
				return
			}
		}
	}
}
