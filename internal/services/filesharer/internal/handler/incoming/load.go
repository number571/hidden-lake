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
	"github.com/number571/hidden-lake/internal/utils/query"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler/incoming/limiters"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

func HandleIncomingLoadHTTP(
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
			_ = api.Response(pW, http.StatusMethodNotAllowed, "method", api.WithoutContentType())
			return
		}

		queryParams := pR.URL.Query()

		fileName := queryParams.Get("name")
		if utils.FileNameIsInvalid(fileName) {
			pLogger.PushWarn(logBuilder.WithMessage("got_invalid_name"))
			_ = api.Response(pW, http.StatusBadRequest, "filename", api.WithoutContentType())
			return
		}

		isPersonal, err := query.GetBoolValueFromQuery(queryParams, "personal")
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("parse_personal"))
			_ = api.Response(pW, http.StatusBadRequest, "personal", api.WithoutContentType())
			return
		}

		chunk, err := strconv.ParseUint(queryParams.Get("chunk"), 10, 64)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_chunk"))
			_ = api.Response(pW, http.StatusBadRequest, "readchunk", api.WithoutContentType())
			return
		}

		aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName)
		if aliasName == "" {
			pLogger.PushErro(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadGateway, "alias", api.WithoutContentType())
			return
		}

		stgPath := utils.GetSharingStoragePath(pPathTo, aliasName, isPersonal)
		fullPath := filepath.Join(stgPath, fileName)

		info, err := fileinfo.NewFileInfo(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				pLogger.PushWarn(logBuilder.WithMessage("file_not_found"))
				_ = api.Response(pW, http.StatusNotFound, "notfound", api.WithoutContentType())
				return
			}
			pLogger.PushErro(logBuilder.WithMessage("get_file_info"))
			_ = api.Response(pW, http.StatusInternalServerError, "fileinfo", api.WithoutContentType())
			return
		}

		chunkSize, err := limiters.GetLimitOnLoadResponseSize(pCtx, pHlkClient)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_chunk_size"))
			_ = api.Response(pW, http.StatusBadGateway, "sizechunk", api.WithoutContentType())
			return
		}

		chunksCount := getChunksCount(info.GetSize(), chunkSize) //nolint:gosec
		if chunk >= chunksCount {
			pLogger.PushWarn(logBuilder.WithMessage("chunk_number"))
			_ = api.Response(pW, http.StatusLengthRequired, "getchunk", api.WithoutContentType())
			return
		}

		file, err := os.Open(fullPath) //nolint:gosec
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("open_file"))
			_ = api.Response(pW, http.StatusInternalServerError, "openfile", api.WithoutContentType())
			return
		}
		defer func() { _ = file.Close() }()

		chunkOffset := int64(chunk) * int64(chunkSize) //nolint:gosec
		nS, err := file.Seek(chunkOffset, io.SeekStart)
		if err != nil || nS != chunkOffset {
			pLogger.PushErro(logBuilder.WithMessage("seek_file"))
			_ = api.Response(pW, http.StatusInternalServerError, "seekfile", api.WithoutContentType())
			return
		}

		buf := make([]byte, chunkSize)
		nR, err := file.Read(buf)
		if err != nil || (chunk != chunksCount-1 && uint64(nR) != chunkSize) { //nolint:gosec
			pLogger.PushErro(logBuilder.WithMessage("chunk_number"))
			_ = api.Response(pW, http.StatusInternalServerError, "readchunk", api.WithoutContentType())
			return
		}

		fileHash := info.GetHash()[:hls_settings.CTrimHashSizeOnLoad]
		pW.Header().Set(hls_settings.CHeaderFileHashOnLoad, fileHash)

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, buf[:nR], api.WithoutContentType())
	}
}

func getChunksCount(pBytesNum, pChunkSize uint64) uint64 {
	return uint64(math.Ceil(float64(pBytesNum) / float64(pChunkSize)))
}
