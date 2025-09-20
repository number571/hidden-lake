package handler

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"

	hls_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
)

func HandleIncomingLoadHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pPathTo string,
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hls_filesharer_settings.CAppShortName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		query := pR.URL.Query()

		name := filepath.Base(query.Get("name"))
		if name != query.Get("name") {
			pLogger.PushWarn(logBuilder.WithMessage("got_another_name"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: got another name")
			return
		}

		chunk, err := strconv.Atoi(query.Get("chunk"))
		if err != nil || chunk < 0 {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_chunk"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: incorrect chunk")
			return
		}

		fullPath := filepath.Join(pPathTo, hls_filesharer_settings.CPathSTG, name)
		stat, err := os.Stat(fullPath)
		if os.IsNotExist(err) || stat.IsDir() {
			pLogger.PushWarn(logBuilder.WithMessage("file_not_found"))
			_ = api.Response(pW, http.StatusNotFound, "failed: file not found")
			return
		}

		chunkSize, err := utils.GetMessageLimit(pCtx, pHlsClient)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_chunk_size"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get chunk size")
			return
		}

		chunks := utils.GetChunksCount(uint64(stat.Size()), chunkSize) //nolint:gosec
		if uint64(chunk) >= chunks {
			pLogger.PushWarn(logBuilder.WithMessage("chunk_number"))
			_ = api.Response(pW, http.StatusLengthRequired, "failed: chunk number")
			return
		}

		file, err := os.Open(fullPath) //nolint:gosec
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("open_file"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: open file")
			return
		}
		defer func() { _ = file.Close() }()

		buf := make([]byte, chunkSize)
		chunkOffset := int64(chunk) * int64(chunkSize) //nolint:gosec

		nS, err := file.Seek(chunkOffset, io.SeekStart)
		if err != nil || nS != chunkOffset {
			pLogger.PushWarn(logBuilder.WithMessage("seek_file"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: seek file")
			return
		}

		nR, err := file.Read(buf)
		if err != nil || (uint64(chunk) != chunks-1 && uint64(nR) != chunkSize) { //nolint:gosec
			pLogger.PushWarn(logBuilder.WithMessage("chunk_number"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: chunk number")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, buf[:nR])
	}
}
