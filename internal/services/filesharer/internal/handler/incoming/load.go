package incoming

import (
	"context"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler/incoming/limiters"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
)

func HandleIncomingLoadHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
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

		fileName := filepath.Base(queryParams.Get("name"))
		if fileName != queryParams.Get("name") {
			pLogger.PushWarn(logBuilder.WithMessage("got_another_name"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: got another name")
			return
		}

		chunk, err := strconv.ParseUint(queryParams.Get("chunk"), 10, 64)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_chunk"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: incorrect chunk")
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
		stat, err := os.Stat(fullPath)
		if os.IsNotExist(err) || stat.IsDir() {
			pLogger.PushWarn(logBuilder.WithMessage("file_not_found"))
			_ = api.Response(pW, http.StatusNotFound, "failed: file not found")
			return
		}

		chunkSize, err := limiters.GetLimitOnLoadResponseSize(pCtx, pHlkClient)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_chunk_size"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get chunk size")
			return
		}

		chunks := getChunksCount(uint64(stat.Size()), chunkSize) //nolint:gosec
		if chunk >= chunks {
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
		if err != nil || (chunk != chunks-1 && uint64(nR) != chunkSize) { //nolint:gosec
			pLogger.PushWarn(logBuilder.WithMessage("chunk_number"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: chunk number")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, buf[:nR])
	}
}

func getChunksCount(pBytesNum, pChunkSize uint64) uint64 {
	return uint64(math.Ceil(float64(pBytesNum) / float64(pChunkSize)))
}
