package app

import (
	"net/http"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/pinger/internal/handler"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
)

func (p *sApp) initExternalServiceHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hls_pinger_settings.CPingPath,
		handler.HandleIncomingPingHTTP(p.fConfig, p.fHTTPLogger),
	) // POST

	buildSettings := build.GetSettings()
	p.fExtServiceHTTP = &http.Server{
		Addr:         p.fConfig.GetAddress().GetExternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpHandleTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpReadTimeout(),
		WriteTimeout: buildSettings.GetHttpHandleTimeout(),
	}
}
