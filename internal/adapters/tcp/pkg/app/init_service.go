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
	cfg := p.fWrapper.GetConfig()

	mux.HandleFunc(
		hla_settings.CHandleIndexPath,
		handler.HandleIndexAPI(p.fHTTPLogger),
	)
	mux.HandleFunc(
		hla_settings.CHandleConfigSettingsPath,
		handler.HandleConfigSettingsAPI(cfg, p.fHTTPLogger),
	)
	mux.HandleFunc(
		hla_settings.CHandleConfigConnectsPath,
		handler.HandleConfigConnectsAPI(pCtx, p.fWrapper, p.fHTTPLogger, p.fTCPAdapter),
	)
	mux.HandleFunc(
		hla_settings.CHandleNetworkOnlinePath,
		handler.HandleNetworkOnlineAPI(p.fHTTPLogger, p.fTCPAdapter),
	)
	mux.HandleFunc(
		hla_settings.CHandleNetworkAdapterPath,
		handler.HandleNetworkProduceAPI(pCtx, cfg, p.fHTTPLogger, p.fTCPAdapter),
	)

	p.fServiceHTTP = &http.Server{
		Addr:        cfg.GetAddress().GetHTTP(),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
