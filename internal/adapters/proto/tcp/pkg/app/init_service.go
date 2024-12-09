package app

import (
	"context"
	"net/http"
	"time"

	hla_settings "github.com/number571/hidden-lake/internal/adapters/pkg/settings"
	"github.com/number571/hidden-lake/internal/adapters/proto/tcp/internal/handler"
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
		Addr:        cfg.GetAddress().GetInternal(),
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
