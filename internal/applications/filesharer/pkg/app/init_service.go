package app

import (
	"context"
	"net/http"
	"os"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/applications/filesharer/internal/handler"
	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app/config"
	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	"github.com/number571/hidden-lake/internal/webui"
)

func (p *sApp) initExternalServiceHTTP(pCtx context.Context, pHlsClient hls_client.IClient) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlf_settings.CLoadPath,
		handler.HandleIncomingLoadHTTP(pCtx, p.fHTTPLogger, p.fStgPath, pHlsClient),
	) // POST

	mux.HandleFunc(
		hlf_settings.CListPath,
		handler.HandleIncomingListHTTP(p.fHTTPLogger, p.fConfig, p.fStgPath),
	) // POST

	buildSettings := build.GetSettings()
	p.fExtServiceHTTP = &http.Server{
		Addr:         p.fConfig.GetAddress().GetExternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpHandleTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpReadTimeout(),
		WriteTimeout: buildSettings.GetHttpWriteTimeout(),
	}
}

func (p *sApp) initInternalServiceHTTP(pCtx context.Context, pHlsClient hls_client.IClient) {
	mux := http.NewServeMux()
	mux.Handle(hlf_settings.CStaticPath, http.StripPrefix(
		hlf_settings.CStaticPath,
		handleFileServer(p.fHTTPLogger, p.fConfig, http.FS(webui.MustGetStaticPath()))),
	)

	cfgWrapper := config.NewWrapper(p.fConfig)

	mux.HandleFunc(hlf_settings.CHandleIndexPath, handler.IndexPage(p.fHTTPLogger, p.fConfig))                              // GET, POST
	mux.HandleFunc(hlf_settings.CHandleAboutPath, handler.AboutPage(p.fHTTPLogger, p.fConfig))                              // GET
	mux.HandleFunc(hlf_settings.CHandleSettingsPath, handler.SettingsPage(pCtx, p.fHTTPLogger, cfgWrapper, pHlsClient))     // GET, PATCH, PUT, POST, DELETE
	mux.HandleFunc(hlf_settings.CHandleFriendsPath, handler.FriendsPage(pCtx, p.fHTTPLogger, p.fConfig, pHlsClient))        // GET, POST, DELETE
	mux.HandleFunc(hlf_settings.CHandleFriendsStoragePath, handler.StoragePage(pCtx, p.fHTTPLogger, p.fConfig, pHlsClient)) // GET, POST, DELETE

	buildSettings := build.GetSettings()
	p.fIntServiceHTTP = &http.Server{
		Addr:         p.fConfig.GetAddress().GetInternal(),
		Handler:      mux, // http.TimeoutHandler returns bug with progress bar of file downloading
		ReadTimeout:  buildSettings.GetHttpReadTimeout(),
		WriteTimeout: buildSettings.GetHttpWriteTimeout(),
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
