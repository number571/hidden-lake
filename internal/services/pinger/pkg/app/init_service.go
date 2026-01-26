package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/pinger/internal/handler"
	"github.com/number571/hidden-lake/internal/services/pinger/internal/handler/incoming"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
)

func (p *sApp) initExternalServiceHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_pinger_settings.CPingPath,
		incoming.HandleIncomingPingHTTP(p.fConfig, p.fHTTPLogger),
	) // GET

	buildSettings := build.GetSettings()
	p.fExtServiceHTTP = &http.Server{
		Addr:         p.fConfig.GetAddress().GetExternal(),
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
		hls_pinger_settings.CHandleIndexPath,
		handler.HandleIndexAPI(p.fHTTPLogger),
	) // GET

	mux.HandleFunc(
		hls_pinger_settings.CHandleCommandPingPath,
		handler.HandleCommandPingAPI(pCtx, p.fConfig, p.fHTTPLogger, pHlkClient),
	) // GET

	buildSettings := build.GetSettings()
	p.fIntServiceHTTP = &http.Server{ // nolint: gosec
		Addr:         p.fConfig.GetAddress().GetInternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpCallbackTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpCallbackTimeout(),
		WriteTimeout: buildSettings.GetHttpCallbackTimeout(),
	}
}
