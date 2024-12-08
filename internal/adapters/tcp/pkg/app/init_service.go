package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/hidden-lake/internal/adapters/tcp/internal/handler"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hla_settings.CHandleNetworkProducePath,
		handler.HandleNetworkProduceAPI(pCtx, p.fWrapper.GetConfig(), p.fHTTPLogger, p.fTCPAdapter),
	)
	mux.HandleFunc(
		hla_settings.CHandleConfigConnectsPath,
		handler.HandleConfigConnectsAPI(pCtx, p.fWrapper, p.fHTTPLogger, p.fTCPAdapter),
	)

	p.fServiceHTTP = &http.Server{
		Addr:        p.fWrapper.GetConfig().GetAddress().GetHTTP(),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
