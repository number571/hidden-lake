package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/utils"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_request "github.com/number571/hidden-lake/pkg/network/request"
)

func HandleStorageFileInfoAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
	pPathTo string,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)
		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		queryParams := pR.URL.Query()

		aliasName := queryParams.Get("friend")
		info, err := getFileInfo(
			pCtx,
			pHlkClient,
			aliasName,
			queryParams.Get("name"),
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_file_info"))
			_ = api.Response(pW, http.StatusForbidden, "failed: get file info")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, info)
	}
}

func getFileInfo(
	pCtx context.Context,
	pHlkClient hlk_client.IClient,
	pAliasName string,
	pFileName string,
) (utils.IFileInfo, error) {
	req := hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_settings.CAppShortName).
		WithPath(fmt.Sprintf("%s?name=%s", hls_settings.CInfoPath, pFileName)).
		Build()

	resp, err := pHlkClient.FetchRequest(pCtx, pAliasName, req)
	if err != nil {
		return nil, err
	}

	if resp.GetCode() != http.StatusOK {
		return nil, errors.New("bad status code") // nolint: err113
	}

	info := &utils.SFileInfo{}
	if err := encoding.DeserializeJSON(resp.GetBody(), info); err != nil {
		return nil, err
	}

	return info, nil
}
