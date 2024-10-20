package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/hidden-lake/internal/service/internal/handler"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()
	cfg := p.fCfgW.GetConfig()

	mux.HandleFunc(hls_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hls_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(p.fCfgW, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(pCtx, p.fCfgW, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(p.fCfgW, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(pCtx, cfg, p.fHTTPLogger, p.fNode))
	mux.HandleFunc(hls_settings.CHandleServicePubKeyPath, handler.HandleServicePubKeyAPI(p.fHTTPLogger, p.fNode))

	p.fServiceHTTP = &http.Server{
		Addr:        cfg.GetAddress().GetHTTP(),
		Handler:     mux,
		ReadTimeout: (5 * time.Second),
	}
}
