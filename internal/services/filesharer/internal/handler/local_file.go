package handler

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/utils"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
)

func HandleLocalFileAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHlkClient hlk_client.IClient,
	pPathTo string,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		queryParams := pR.URL.Query()
		aliasName := queryParams.Get("friend")

		fileName := filepath.Base(queryParams.Get("name"))
		if fileName != queryParams.Get("name") {
			pLogger.PushWarn(logBuilder.WithMessage("got_another_name"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: got another name")
			return
		}

		stgPath, err := utils.GetSharingStoragePath(pCtx, pPathTo, pHlkClient, aliasName, aliasName != "")
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("get_path_to_file"))
			_ = api.Response(pW, http.StatusForbidden, "failed: get path to file")
			return
		}

		fullPath := filepath.Join(stgPath, fileName)

		switch pR.Method {
		case http.MethodGet:
			file, err := os.Open(fullPath) // nolint: gosec
			if err != nil {
				pLogger.PushErro(logBuilder.WithMessage("read_file"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: read file")
				return
			}
			defer func() { _ = file.Close() }()

			if err := api.ResponseWithReader(pW, http.StatusOK, file); err != nil {
				pLogger.PushErro(logBuilder.WithMessage("stream_reader"))
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))

		case http.MethodDelete:
			if err := os.Remove(fullPath); err != nil {
				pLogger.PushErro(logBuilder.WithMessage("delete_file"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: delete file")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: delete file")

		case http.MethodPost:
			if err := os.MkdirAll(stgPath, 0700); err != nil {
				pLogger.PushErro(logBuilder.WithMessage("mkdir_all"))
				_ = api.Response(pW, http.StatusForbidden, "failed: mkdir all")
				return
			}

			dst, err := os.Create(fullPath) // nolint: gosec
			if err != nil {
				pLogger.PushErro(logBuilder.WithMessage("create_file"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: create file")
				return
			}
			defer func() { _ = dst.Close() }()

			if _, err := io.Copy(dst, pR.Body); err != nil {
				pLogger.PushErro(logBuilder.WithMessage("copy_file"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: copy file")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: upload file")
		}
	}
}
