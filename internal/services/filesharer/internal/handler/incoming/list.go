package incoming

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
)

func HandleIncomingListHTTP(
	pLogger logger.ILogger,
	pCfg config.IConfig,
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

		page, err := strconv.Atoi(pR.URL.Query().Get("page"))
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_page"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: incorrect page")
			return
		}

		result, err := getListFileInfo(pCfg, pPathTo, uint64(page)) //nolint:gosec
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("open storage"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: open storage")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, result)
	}
}

func getListFileInfo(pCfg config.IConfig, pPathTo string, pPage uint64) ([]fileinfo.IFileInfo, error) {
	pageOffset := pCfg.GetSettings().GetPageOffset()
	fileReader := pageOffset

	stgPath := filepath.Join(pPathTo, hls_filesharer_settings.CPathSTG)
	entries, err := os.ReadDir(stgPath)
	if err != nil {
		return nil, err
	}

	files := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, e)
	}

	result := make([]fileinfo.IFileInfo, 0, pageOffset)
	for i := (pPage * pageOffset); i < uint64(len(files)); i++ {
		if fileReader == 0 {
			break
		}
		fileReader--

		fileName := files[i].Name()
		fullPath := filepath.Join(stgPath, fileName)

		info, err := fileinfo.NewFileInfo(fullPath)
		if err != nil {
			return nil, err
		}

		result = append(result, info)
	}

	return result, nil
}
