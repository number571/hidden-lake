package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	hls_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/chars"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	"github.com/number571/hidden-lake/pkg/api/services/pinger/request"
)

func HandleCommandPingAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		req := request.NewPingRequest()
		resp, err := pHlkClient.FetchRequest(pCtx, pR.URL.Query().Get("friend"), req)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("fetch_request"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: fetch request")
			return
		}

		if code := resp.GetCode(); code != http.StatusOK {
			pLogger.PushErro(logBuilder.WithMessage("status_error"))
			_ = api.Response(pW, http.StatusTeapot, fmt.Sprintf("failed: status %d", code))
			return
		}

		respBody := string(resp.GetBody())
		if chars.HasNotGraphicCharacters(respBody) {
			pLogger.PushErro(logBuilder.WithMessage("invalid_response"))
			_ = api.Response(pW, http.StatusServiceUnavailable, "failed: invalid response")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, respBody)
	}
}
