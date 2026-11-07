package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/publisher/internal/database"
	"github.com/number571/hidden-lake/internal/services/publisher/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/query"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	"github.com/number571/hidden-lake/pkg/api/services/publisher/request"
)

func HandleMsgboxSizeAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pHlkClient hlk_client.IClient,
	pDatabase database.IKVDatabase,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		queryParams := pR.URL.Query()

		aliasName := queryParams.Get("friend")
		if aliasName == "" {
			pLogger.PushWarn(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: alias name is null")
			return
		}

		isRemote, err := query.GetBoolValueFromQuery(queryParams, "remote")
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("remote_error"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: remote is unknown")
			return
		}

		if !isRemote {
			size := pDatabase.Size(aliasName)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, strconv.FormatUint(size, 10))
			return
		}

		req := request.NewMsgboxSizeRequest()
		rsp, err := pHlkClient.FetchRequest(pCtx, aliasName, req)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("no_response"))
			_ = api.Response(pW, http.StatusServiceUnavailable, "failed: no response")
			return
		}
		if rsp.GetCode() != http.StatusOK {
			pLogger.PushWarn(logBuilder.WithMessage("bad_status_code"))
			_ = api.Response(pW, http.StatusTeapot, "failed: status code")
			return
		}
		size, err := strconv.ParseUint(string(rsp.GetBody()), 10, 64)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("bad_response"))
			_ = api.Response(pW, http.StatusTeapot, "failed: body response")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, strconv.FormatUint(size, 10))
	}
}
