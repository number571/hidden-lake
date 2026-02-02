package incoming

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"

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

		logBuilder := http_logger.NewLogBuilder(hls_filesharer_settings.GetAppShortNameFMT(), pR)

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

		page, err := strconv.Atoi(queryParams.Get("page"))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_page"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: incorrect page")
			return
		}

		aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName)
		stgPath, err := utils.GetSharingStoragePath(pCtx, pPathTo, pHlkClient, aliasName, isPersonal)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_path_to_file"))
			_ = api.Response(pW, http.StatusForbidden, "failed: get path to file")
			return
		}

		stat, err := os.Stat(stgPath)
		if os.IsNotExist(err) || !stat.IsDir() {
			list, err := dto.LoadFileInfoList("[]")
			if err != nil {
				panic(err)
			}
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, list.ToString())
			return
		}

		list, err := utils.GetFileInfoList(stgPath, uint64(page), pCfg.GetSettings().GetPageOffset()) //nolint:gosec
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("open storage"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: open storage")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, list.ToString())
	}
}
