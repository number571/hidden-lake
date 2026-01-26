package incoming

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

func HandleIncomingInfoHTTP(
	pLogger logger.ILogger,
	pPathTo string,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hls_filesharer_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		query := pR.URL.Query()

		fileName := filepath.Base(query.Get("name"))
		if fileName != query.Get("name") {
			pLogger.PushWarn(logBuilder.WithMessage("got_another_name"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: got another name")
			return
		}

		stgPath := filepath.Join(pPathTo, hls_filesharer_settings.CPathSTG)
		fullPath := filepath.Join(stgPath, fileName)

		stat, err := os.Stat(fullPath)
		if os.IsNotExist(err) || stat.IsDir() {
			pLogger.PushWarn(logBuilder.WithMessage("file_not_found"))
			_ = api.Response(pW, http.StatusNotFound, "failed: file not found")
			return
		}

		info, err := fileinfo.NewFileInfo(fullPath)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_file_info"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: get file info")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, info)
	}
}
