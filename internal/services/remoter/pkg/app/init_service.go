package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/remoter/internal/handler"
	"github.com/number571/hidden-lake/internal/services/remoter/internal/handler/incoming"
	hls_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"
)

func (p *sApp) initExternalServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_settings.CExecPath,
		incoming.HandleIncomingExecHTTP(pCtx, p.fConfig, p.fHTTPLogger),
	) // POST

	buildSettings := build.GetSettings()
	p.fExtServiceHTTP = &http.Server{
		Addr: p.fConfig.GetAddress().GetExternal(),
		// no need http_handle_timeout -> used custom exec_timeout
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpHandleTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpReadTimeout(),
		WriteTimeout: buildSettings.GetHttpHandleTimeout(),
	}
}

func (p *sApp) initInternalServiceHTTP(
	pCtx context.Context,
	pHlkClient hlk_client.IClient,
) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_settings.CHandleIndexPath,
		handler.HandleIndexAPI(p.fHTTPLogger),
	) // GET

	mux.HandleFunc(
		hls_settings.CHandleCommandExecPath,
		handler.HandleCommandExecAPI(pCtx, p.fConfig, p.fHTTPLogger, pHlkClient),
	) // POST

	buildSettings := build.GetSettings()
	p.fIntServiceHTTP = &http.Server{ // nolint: gosec
		Addr:         p.fConfig.GetAddress().GetInternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpCallbackTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpCallbackTimeout(),
		WriteTimeout: buildSettings.GetHttpCallbackTimeout(),
	}
}
