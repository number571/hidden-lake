package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/client/fileinfo"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func HandleStorageListAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)
		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		page, err := strconv.ParseUint(pR.URL.Query().Get("page"), 10, 64)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("parse_page"))
			_ = api.Response(pW, http.StatusForbidden, "failed: parse page")
			return
		}

		req := hlk_request.NewRequestBuilder().
			WithMethod(http.MethodGet).
			WithHost(hls_settings.CAppShortName).
			WithPath(fmt.Sprintf("%s?page=%d", hls_settings.CListPath, page)).
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

		infos, err := fileinfo.LoadFileInfoList(resp.GetBody())
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("decode_response"))
			_ = api.Response(pW, http.StatusForbidden, "failed: decode response")
			return
		}

		_ = api.Response(pW, http.StatusOK, infos)
	}
}
