package app

import (
	"net/http"
	"time"

	"github.com/number571/hidden-lake/internal/applications/pinger/internal/handler"
	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"
)

func (p *sApp) initExternalServiceHTTP() {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlp_settings.CPingPath,
		handler.HandleIncomingPingHTTP(p.fConfig, p.fHTTPLogger),
	) // POST

	p.fExtServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetExternal(),
		Handler:     http.TimeoutHandler(mux, time.Minute, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
