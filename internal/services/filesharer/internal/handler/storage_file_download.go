package handler

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/number571/go-peer/pkg/logger"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler/stream"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/client/fileinfo"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/pubkey"
)

func HandleStorageFileDownloadAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
	pPathTo string,
) http.HandlerFunc {
	downloadProcessesMap := newDownloadProcessesMap()

	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		queryParams := pR.URL.Query()
		aliasName := queryParams.Get("friend")

		req := newFileInfoRequest(queryParams.Get("name"))
		resp, err := pHlkClient.FetchRequest(pCtx, aliasName, req)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("fetch_request"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: fetch request")
			return
		}

		if resp.GetCode() != http.StatusOK {
			pLogger.PushErro(logBuilder.WithMessage("status_error"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: status error")
			return
		}

		info, err := fileinfo.LoadFileInfo(resp.GetBody())
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("decode_response"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: decode response")
			return
		}

		fileHash := info.GetHash()
		pW.Header().Set(hls_settings.CHeaderFileHash, fileHash)

		if ok := downloadProcessesMap.Exist(fileHash); ok {
			pW.Header().Set(hls_settings.CHeaderInProcess, hls_settings.CHeaderProcessModeY)
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusAccepted, "process: download")
			return
		}

		pubKey, err := pubkey.GetFriendPubKeyByAliasName(pCtx, pHlkClient, aliasName)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_pubkey"))
			_ = api.Response(pW, http.StatusForbidden, "failed: get public key")
			return
		}

		pathToDownload := filepath.Join(pPathTo, pubKey.GetHasher().ToString())
		if err := os.MkdirAll(pathToDownload, 0700); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("mkdir_all"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: mkdir all")
			return
		}

		reader, err := stream.BuildStreamReader(
			pCtx,
			pConfig.GetSettings().GetRetryNum(),
			pathToDownload,
			aliasName,
			pHlkClient,
			info,
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("build_stream"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: build stream")
			return
		}

		downloadProcessesMap.Set(fileHash)
		defer downloadProcessesMap.Del(fileHash)

		pW.Header().Set(hls_settings.CHeaderInProcess, hls_settings.CHeaderProcessModeN)
		if err := api.ResponseWithReader(pW, http.StatusOK, reader); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("stream_reader"))
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
	}
}

type downloadProcessesMap struct {
	fMutex *sync.RWMutex
	fMap   map[string]struct{}
}

func newDownloadProcessesMap() *downloadProcessesMap {
	return &downloadProcessesMap{
		fMutex: &sync.RWMutex{},
		fMap:   make(map[string]struct{}, 256),
	}
}

func (p *downloadProcessesMap) Set(k string) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	p.fMap[k] = struct{}{}
}

func (p *downloadProcessesMap) Del(k string) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	delete(p.fMap, k)
}

func (p *downloadProcessesMap) Exist(k string) bool {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	_, ok := p.fMap[k]
	return ok
}
