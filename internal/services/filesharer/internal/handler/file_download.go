package handler

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/logger"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler/stream"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/pubkey"
)

func HandleFileDownloadAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
	pPathTo string,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)
		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		aliasName := pR.URL.Query().Get("friend")
		pubKey, err := pubkey.GetFriendPubKeyByAliasName(pCtx, pHlkClient, aliasName)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_pubkey"))
			_ = api.Response(pW, http.StatusForbidden, "failed: get public key")
			return
		}

		pathToDownload := filepath.Join(pPathTo, pubKey.GetHasher().ToString())
		if err := os.MkdirAll(pathToDownload, 0700); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("mkdir_all"))
			_ = api.Response(pW, http.StatusForbidden, "failed: mkdir all")
			return
		}

		info, err := getFileInfo(
			pCtx,
			pHlkClient,
			pR.URL.Query().Get("name"),
			aliasName,
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_file_info"))
			_ = api.Response(pW, http.StatusForbidden, "failed: get file info")
			return
		}

		reader, tempFIle, err := stream.BuildStreamReader(
			pCtx,
			pConfig.GetSettings().GetRetryNum(),
			pathToDownload,
			aliasName,
			pHlkClient,
			info,
		)
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("build_stream"))
			_ = api.Response(pW, http.StatusForbidden, "failed: build stream")
			return
		}
		defer func() { _ = os.Remove(tempFIle) }()

		if _, err := io.Copy(pW, reader); err != nil {
			pLogger.PushErro(logBuilder.WithMessage("copy_file"))
			_ = api.Response(pW, http.StatusForbidden, "failed: copy file")
			return
		}
	}
}
