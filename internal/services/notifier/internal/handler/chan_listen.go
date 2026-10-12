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

		aliasName := queryParams.Get("friend")
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: alias name is null")
			return
		}

		sid := queryParams.Get("sid")
		if err := pBroker.Register(sid); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("limit_subs"))
			_ = api.Response(pW, http.StatusNotAcceptable, "has limit of subscribers")
			return
		}

		ctx, cancel := context.WithTimeout(pCtx, buildSettings.GetHttpReadTimeout())
		defer cancel()

		for {
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

			if mc.GetFriend() != aliasName {
				pLogger.PushInfo(logBuilder.WithMessage("another_friend"))
				continue
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, mc.GetMessage().ToString(), api.WithCustomContentType(api.CApplicationJSON))
			return
		}
	}
}
