package app

import (
	"context"
	"net/http"
	"os"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	hls_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/app/config"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/webui"
)

func (p *sApp) initExternalServiceHTTP(pCtx context.Context, pHlsClient hls_client.IClient) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hls_filesharer_settings.CLoadPath,
		handler.HandleIncomingLoadHTTP(pCtx, p.fHTTPLogger, p.fPathTo, pHlsClient),
	) // POST

	mux.HandleFunc(
		hls_filesharer_settings.CListPath,
		handler.HandleIncomingListHTTP(p.fHTTPLogger, p.fConfig, p.fPathTo),
	) // POST

	buildSettings := build.GetSettings()
	p.fExtServiceHTTP = &http.Server{
		Addr:         p.fConfig.GetAddress().GetExternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpHandleTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpReadTimeout(),
		WriteTimeout: buildSettings.GetHttpHandleTimeout(),
	}
}

func (p *sApp) initInternalServiceHTTP(pCtx context.Context, pHlsClient hls_client.IClient) {
	mux := http.NewServeMux()
	mux.Handle(hls_filesharer_settings.CStaticPath, http.StripPrefix(
		hls_filesharer_settings.CStaticPath,
		handleFileServer(p.fHTTPLogger, p.fConfig, http.FS(webui.MustGetStaticPath()))),
	)

	cfgWrapper := config.NewWrapper(p.fConfig)

	mux.HandleFunc(hls_filesharer_settings.CHandleIndexPath, handler.IndexPage(p.fHTTPLogger, p.fConfig))                                         // GET, POST
	mux.HandleFunc(hls_filesharer_settings.CHandleAboutPath, handler.AboutPage(p.fHTTPLogger, p.fConfig))                                         // GET
	mux.HandleFunc(hls_filesharer_settings.CHandleSettingsPath, handler.SettingsPage(pCtx, p.fHTTPLogger, cfgWrapper, pHlsClient))                // GET, PATCH, PUT, POST, DELETE
	mux.HandleFunc(hls_filesharer_settings.CHandleFriendsPath, handler.FriendsPage(pCtx, p.fHTTPLogger, p.fConfig, pHlsClient))                   // GET, POST, DELETE
	mux.HandleFunc(hls_filesharer_settings.CHandleFriendsStoragePath, handler.StoragePage(pCtx, p.fHTTPLogger, p.fConfig, p.fPathTo, pHlsClient)) // GET, POST, DELETE

	buildSettings := build.GetSettings()
	p.fIntServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetInternal(),
		Handler:     mux,
		ReadTimeout: buildSettings.GetHttpReadTimeout(),
		// WriteTimeout not set (downloading a file may take a long time)
	}
}

func handleFileServer(pLogger logger.ILogger, pCfg config.IConfig, pFS http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := pFS.Open(r.URL.Path); os.IsNotExist(err) {
			handler.NotFoundPage(pLogger, pCfg)(w, r)
			return
		}
		http.FileServer(pFS).ServeHTTP(w, r)
	})
}
