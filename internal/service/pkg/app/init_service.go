package app

import (
	"context"
	"net/http"
	"time"

	"github.com/number571/hidden-lake/internal/adapters/tcp/pkg/client"
	"github.com/number571/hidden-lake/internal/service/internal/handler"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	mux := http.NewServeMux()
	cfg := p.fCfgW.GetConfig()
	origNode := p.fNode.GetAnonymityNode()

	var hlaClient client.IClient
	for _, a := range cfg.GetAdapters() {
		hlaClient = client.NewClient(
			client.NewRequester(
				a, // TODO:
				&http.Client{Timeout: 5 * time.Second},
			),
		)
		break
	}

	mux.HandleFunc(hls_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hls_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(p.fCfgW, p.fHTTPLogger, origNode))
	mux.HandleFunc(hls_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(pCtx, p.fHTTPLogger, hlaClient))
	mux.HandleFunc(hls_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(p.fCfgW, p.fHTTPLogger, origNode))
	mux.HandleFunc(hls_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(p.fHTTPLogger, origNode))
	mux.HandleFunc(hls_settings.CHandleServicePubKeyPath, handler.HandleServicePubKeyAPI(p.fHTTPLogger, origNode))
	mux.HandleFunc(hls_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(pCtx, cfg, p.fHTTPLogger, p.fNode))

	p.fServiceHTTP = &http.Server{
		Addr:        cfg.GetAddress().GetInternal(),
		Handler:     mux,
		ReadTimeout: (5 * time.Second),
	}
}
