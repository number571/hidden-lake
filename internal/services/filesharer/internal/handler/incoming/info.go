package incoming

import (
	"crypto/sha512"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/utils"
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

		info := utils.NewFileInfo(
			fileName,
			getFileHash(fullPath),
			getFileSize(fullPath),
		)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, info)
	}
}

func getFileSize(filename string) uint64 {
	stat, _ := os.Stat(filename)
	return uint64(stat.Size()) //nolint:gosec
}

func getFileHash(filename string) string {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		return ""
	}
	defer func() { _ = f.Close() }()

	h := sha512.New384()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return encoding.HexEncode(h.Sum(nil))
}
