package app

import (
	"net/http"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/applications/pinger/internal/handler"
	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
)

func (p *sApp) initExternalServiceHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlp_settings.CPingPath,
		handler.HandleIncomingPingHTTP(p.fConfig, p.fHTTPLogger),
	) // POST

	buildSettings := build.GetSettings()
	p.fExtServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetExternal(),
		Handler:     http.TimeoutHandler(mux, buildSettings.GetHttpHandleTimeout(), "handle timeout"),
		ReadTimeout: buildSettings.GetHttpReadTimeout(),
	}
}
