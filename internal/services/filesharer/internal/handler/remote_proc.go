package handler

import (
	"context"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler/process"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

func HandleRemoteFileProcAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pProcessManager process.IDownloadProcessManager,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		queryParams := pR.URL.Query()

		fileName := queryParams.Get("name")
		aliasName := queryParams.Get("friend")

		isPersonal, err := utils.GetBoolValueFromQuery(queryParams, "personal")
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("parse_personal"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: parse personal")
			return
		}

		dpKey := dto.NewDownloadProcessKey(aliasName, fileName, isPersonal)

		if pR.Method == http.MethodDelete {
			if ok := pProcessManager.Unlock(dpKey); !ok {
				pLogger.PushErro(logBuilder.WithMessage("cancel_download"))
				_ = api.Response(pW, http.StatusNotFound, "failed: cancel download")
				return
			}
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: cancel download")
		}

		val, ok := pProcessManager.Get(dpKey)
		if !ok {
			pLogger.PushWarn(logBuilder.WithMessage("load_process"))
			_ = api.Response(pW, http.StatusNotFound, "failed: load process")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, val.ToString())
	}
}
