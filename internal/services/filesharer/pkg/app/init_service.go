package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler/incoming"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
)

func (p *sApp) initExternalServiceHTTP(pCtx context.Context, pHlkClient hlk_client.IClient) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_filesharer_settings.CLoadPath,
		incoming.HandleIncomingLoadHTTP(pCtx, p.fHTTPLogger, p.fPathTo, pHlkClient),
	) // GET

	mux.HandleFunc(
		hls_filesharer_settings.CListPath,
		incoming.HandleIncomingListHTTP(p.fHTTPLogger, p.fConfig, p.fPathTo),
	) // GET

	mux.HandleFunc(
		hls_filesharer_settings.CInfoPath,
		incoming.HandleIncomingInfoHTTP(p.fHTTPLogger, p.fPathTo),
	) // GET

	buildSettings := build.GetSettings()
	p.fExtServiceHTTP = &http.Server{
		Addr:         p.fConfig.GetAddress().GetExternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpHandleTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpReadTimeout(),
		WriteTimeout: buildSettings.GetHttpHandleTimeout(),
	}
}

func (p *sApp) initInternalServiceHTTP(pCtx context.Context, pHlkClient hlk_client.IClient) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_filesharer_settings.CHandleIndexPath,
		handler.HandleIndexAPI(p.fHTTPLogger),
	) // GET

	mux.HandleFunc(
		hls_filesharer_settings.CHandleFileInfoPath,
		handler.HandleFileInfoAPI(pCtx, p.fConfig, p.fHTTPLogger, pHlkClient),
	) // GET

	mux.HandleFunc(
		hls_filesharer_settings.CHandleStorageFilesPath,
		handler.HandleStorageFilesAPI(pCtx, p.fConfig, p.fHTTPLogger, pHlkClient),
	) // GET

	mux.HandleFunc(
		hls_filesharer_settings.CHandleFileDownloadPath,
		handler.HandleFileDownloadAPI(pCtx, p.fConfig, p.fHTTPLogger, pHlkClient, p.fPathTo),
	) // GET

	p.fIntServiceHTTP = &http.Server{ // nolint: gosec
		Addr:    p.fConfig.GetAddress().GetInternal(),
		Handler: mux,
	}
}
