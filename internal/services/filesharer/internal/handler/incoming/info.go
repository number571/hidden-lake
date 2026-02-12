package incoming

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

func HandleIncomingInfoHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
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
		fileName := queryParams.Get("name")

		if utils.FileNameIsInvalid(fileName) {
			pLogger.PushWarn(logBuilder.WithMessage("got_invalid_name"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: got invalid name")
			return
		}

		isPersonal, err := utils.GetBoolValueFromQuery(queryParams, "personal")
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("parse_personal"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: parse personal")
			return
		}

		aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName)
		stgPath, err := utils.GetSharingStoragePath(pCtx, pPathTo, pHlkClient, aliasName, isPersonal)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_path_to_file"))
			_ = api.Response(pW, http.StatusForbidden, "failed: get path to file")
			return
		}

		fullPath := filepath.Join(stgPath, fileName)
		info, err := fileinfo.NewFileInfo(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				pLogger.PushWarn(logBuilder.WithMessage("file_not_found"))
				_ = api.Response(pW, http.StatusNotFound, "failed: file not found")
				return
			}
			pLogger.PushWarn(logBuilder.WithMessage("get_file_info"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: get file info")
			return
		}

		if info.GetName() != fileName {
			pLogger.PushErro(logBuilder.WithMessage("invalid_response"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: invalid response")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, info.ToString())
	}
}
