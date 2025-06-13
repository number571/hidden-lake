package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/applications/remoter/internal/handler"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
)

func (p *sApp) initExternalServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlr_settings.CExecPath,
		handler.HandleIncomingExecHTTP(pCtx, p.fConfig, p.fHTTPLogger),
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
