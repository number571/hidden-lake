package handler

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/query"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/limiter"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/request"
)

func HandleChanNotifyAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			size, err := limiter.GetLimitOnNotifyRequestSize(pCtx, pHlkClient)
			if err != nil {
				pLogger.PushErro(logBuilder.WithMessage("get_limit"))
				_ = api.Response(pW, http.StatusBadGateway, "failed: get limit")
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

		queryParams := pR.URL.Query()

		aliasName := queryParams.Get("friend")
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: alias name is null")
			return
		}

		wait, err := query.GetBoolValueFromQuery(queryParams, "wait")
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("parse_wait"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: parse wait")
			return
		}

		req := request.NewNotifyRequest(body)
		if wait {
			_, err = pHlkClient.FetchRequest(pCtx, aliasName, req)
		} else {
			err = pHlkClient.SendRequest(pCtx, aliasName, req)
		}
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("do_request"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: do request")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, "ok")
	}
}
