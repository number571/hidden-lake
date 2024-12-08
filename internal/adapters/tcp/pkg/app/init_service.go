package app

import (
	"net/http"
	"time"

	"github.com/number571/hidden-lake/internal/adapters/tcp/internal/handler"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
)

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hla_settings.CHandleIndexPath,
		handler.HandleIndexAPI(p.fConfig, p.fHTTPLogger, p.fTCPAdapter),
	)

	p.fServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetHTTP(),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
