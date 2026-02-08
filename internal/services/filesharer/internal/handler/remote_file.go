package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler/stream"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
	"github.com/number571/hidden-lake/pkg/api/services/filesharer/request"
)

func HandleRemoteFileAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
	pPathTo string,
) http.HandlerFunc {
	downloadProcessesMap := newDownloadProcessesMap()

	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		queryParams := pR.URL.Query()

		fileName := queryParams.Get("name")
		aliasName := queryParams.Get("friend")

		isPersonal, err := utils.GetBoolValueFromQuery(queryParams, "personal")
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("parse_personal"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: parse personal")
			return
		}

		stgPath, err := utils.GetPrivateStoragePath(pCtx, pPathTo, pHlkClient, aliasName)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_path_to_file"))
			_ = api.Response(pW, http.StatusForbidden, "failed: get path to file")
			return
		}

		if err := os.MkdirAll(stgPath, 0700); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("mkdir_all"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: mkdir all")
			return
		}

		fullPath := filepath.Join(stgPath, fmt.Sprintf("%s.p%t", fileName, isPersonal))

		if pR.Method == http.MethodDelete {
			if err := os.Remove(fullPath); err != nil {
				pLogger.PushErro(logBuilder.WithMessage("delete_file"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: delete file")
				return
			}
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: delete file")
			return
		}

		req := request.NewInfoRequest(fileName, isPersonal)
		resp, err := pHlkClient.FetchRequest(pCtx, aliasName, req)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("fetch_request"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: fetch request")
			return
		}

		if code := resp.GetCode(); code != http.StatusOK {
			pLogger.PushErro(logBuilder.WithMessage("status_error"))
			_ = api.Response(pW, http.StatusTeapot, fmt.Sprintf("failed: status %d", code))
			return
		}

		info, err := fileinfo.LoadFileInfo(resp.GetBody())
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("decode_response"))
			_ = api.Response(pW, http.StatusTeapot, "failed: decode response")
			return
		}

		if info.GetName() != fileName {
			pLogger.PushErro(logBuilder.WithMessage("invalid_response"))
			_ = api.Response(pW, http.StatusTeapot, "failed: invalid response")
			return
		}

		fileHash := info.GetHash()
		pW.Header().Set(hls_settings.CHeaderFileHash, fileHash)

		if ok := downloadProcessesMap.TryLock(fileHash); !ok {
			pW.Header().Set(hls_settings.CHeaderInProcess, hls_settings.CHeaderProcessModeY)
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusAccepted, "process: download")
			return
		}
		defer downloadProcessesMap.Unlock(fileHash)

		streamReader, err := stream.BuildStreamReader(
			pCtx,
			pConfig.GetSettings().GetRetryNum(),
			fullPath,
			aliasName,
			pHlkClient,
			info,
			isPersonal,
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("build_stream"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: build stream")
			return
		}

		// temp file used as buffer
		if _, err := io.Copy(io.Discard, streamReader); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("read_stream"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: read stream")
			return
		}

		file, err := os.Open(fullPath) // nolint: gosec
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("read_file"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: read file")
			return
		}
		defer func() { _ = file.Close() }()

		pW.Header().Set(hls_settings.CHeaderInProcess, hls_settings.CHeaderProcessModeN)

		if err := api.ResponseWithReader(pW, http.StatusOK, file); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("stream_reader"))
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
	}
}

type downloadProcessesMap struct {
	fMutex *sync.RWMutex
	fMap   map[string]*sync.Mutex
}

func newDownloadProcessesMap() *downloadProcessesMap {
	return &downloadProcessesMap{
		fMutex: &sync.RWMutex{},
		fMap:   make(map[string]*sync.Mutex, 256),
	}
}

func (p *downloadProcessesMap) Unlock(k string) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	mtx, ok := p.fMap[k]
	if !ok {
		panic("unlock mutex without lock")
	}

	mtx.Unlock()
	delete(p.fMap, k)
}

func (p *downloadProcessesMap) TryLock(k string) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	_, ok := p.fMap[k]
	if ok {
		return false
	}

	mtx := &sync.Mutex{}
	p.fMap[k] = mtx
	mtx.Lock()

	return true
}
