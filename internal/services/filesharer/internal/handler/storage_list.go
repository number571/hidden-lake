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

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		page, err := strconv.ParseUint(pR.URL.Query().Get("page"), 10, 64)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("parse_page"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: parse page")
			return
		}

		req := newListFilesRequest(page)
		resp, err := pHlkClient.FetchRequest(pCtx, pR.URL.Query().Get("friend"), req)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("fetch_request"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: fetch request")
			return
		}

		if resp.GetCode() != http.StatusOK {
			pLogger.PushErro(logBuilder.WithMessage("status_error"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: status error")
			return
		}

		infos, err := fileinfo.LoadFileInfoList(resp.GetBody())
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("decode_response"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: decode response")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, infos)
	}
}

func newListFilesRequest(pPage uint64) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf("%s?page=%d", hls_settings.CListPath, pPage)).
		Build()
}
