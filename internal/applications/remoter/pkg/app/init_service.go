package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/hidden-lake/internal/applications/remoter/internal/handler"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
)

func (p *sApp) initIncomingServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hlr_settings.CExecPath,
		handler.HandleIncomingExecHTTP(pCtx, p.fConfig, p.fHTTPLogger),
	) // POST

	execTimeout := p.fConfig.GetSettings().GetExecTimeout()
	p.fIncServiceHTTP = &http.Server{
		Addr:        p.fConfig.GetAddress().GetIncoming(),
		Handler:     http.TimeoutHandler(mux, 2*execTimeout, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
