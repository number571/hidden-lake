package incoming

import (
	"context"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
)

func HandleIncomingListHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pPathTo string,
	pHlkClient hlk_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		queryParams := pR.URL.Query()
		isPersonal, err := utils.GetBoolValueFromQuery(queryParams, "personal")
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("parse_personal"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: parse personal")
			return
		}

		page := uint64(0)
		if v, ok := queryParams["page"]; ok && len(v) > 0 {
			var err error
			page, err = strconv.ParseUint(v[0], 10, 64)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("incorrect_page"))
				_ = api.Response(pW, http.StatusBadRequest, "failed: incorrect page")
				return
			}
		}

		aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName)
		stgPath, err := utils.GetSharingStoragePath(pCtx, pPathTo, pHlkClient, aliasName, isPersonal)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_path_to_file"))
			_ = api.Response(pW, http.StatusForbidden, "failed: get path to file")
			return
		}

		list, err := utils.GetFileInfoList(stgPath, page, pCfg.GetSettings().GetPageOffset()) //nolint:gosec
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("open storage"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: open storage")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, list.ToString())
	}
}
