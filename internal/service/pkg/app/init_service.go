package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/service/internal/handler"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/pkg/adapters/http/client"
)

func (p *sApp) initServiceHTTP(pCtx context.Context) {
	buildSettings := build.GetSettings()

	mux := http.NewServeMux()
	cfg := p.fCfgW.GetConfig()
	origNode := p.fNode.GetOriginNode()

	endpoints := cfg.GetEndpoints()
	epClients := make([]client.IClient, 0, len(endpoints))
	for _, ep := range endpoints {
		requester := client.NewRequester(ep, &http.Client{Timeout: buildSettings.GetHttpHandleTimeout()})
		epClients = append(epClients, client.NewClient(requester))
	}

	mux.HandleFunc(hls_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hls_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(p.fCfgW, p.fHTTPLogger, origNode))
	mux.HandleFunc(hls_settings.CHandleConfigConnectsPath, handler.HandleConfigConnectsAPI(pCtx, p.fHTTPLogger, epClients))
	mux.HandleFunc(hls_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(p.fCfgW, p.fHTTPLogger, origNode))
	mux.HandleFunc(hls_settings.CHandleNetworkOnlinePath, handler.HandleNetworkOnlineAPI(pCtx, p.fHTTPLogger, epClients))
	mux.HandleFunc(hls_settings.CHandleServicePubKeyPath, handler.HandleServicePubKeyAPI(p.fHTTPLogger, origNode))
	mux.HandleFunc(hls_settings.CHandleNetworkRequestPath, handler.HandleNetworkRequestAPI(pCtx, cfg, p.fHTTPLogger, p.fNode))

	// response can take a long time to complete (x2 time QB-problem period)
	callbackTimeout := buildSettings.GetHttpCallbackTimeout()

	p.fServiceHTTP = &http.Server{
		Addr:         cfg.GetAddress().GetInternal(),
		Handler:      http.TimeoutHandler(mux, callbackTimeout, "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpReadTimeout(),
		WriteTimeout: callbackTimeout,
	}
}
