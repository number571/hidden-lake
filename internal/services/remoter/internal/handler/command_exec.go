package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/number571/go-peer/pkg/logger"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/remoter/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_request "github.com/number571/hidden-lake/pkg/request"
)

func HandleCommandExecAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		var vRequest hls_settings.SCommandExecRequest

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vRequest); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		req := hlk_request.NewRequestBuilder().
			WithMethod(http.MethodPost).
			WithHost(hls_settings.CAppShortName).
			WithPath(hls_settings.CExecPath).
			WithHead(map[string]string{
				hls_settings.CHeaderPassword: vRequest.FPassword,
			}).
			WithBody([]byte(strings.Join(vRequest.FCommand, hls_settings.CExecSeparator))).
			Build()

		resp, err := pHlkClient.FetchRequest(pCtx, pR.URL.Query().Get("friend"), req)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("fetch_request"))
			_ = api.Response(pW, http.StatusForbidden, "failed: fetch request")
			return
		}

		if resp.GetCode() != http.StatusOK {
			pLogger.PushErro(logBuilder.WithMessage("status_error"))
			_ = api.Response(pW, http.StatusForbidden, "failed: status error")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, resp.GetBody())
	}
}
