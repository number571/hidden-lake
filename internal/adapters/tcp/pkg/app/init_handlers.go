package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/internal/adapters/tcp/internal/handler"
	hla_settings "github.com/number571/hidden-lake/pkg/adapters/http/settings"
)

func (p *sApp) initHandlers(pCtx context.Context) {
	networkNode := p.fTCPAdapter.GetConnKeeper().GetNetworkNode()
	p.fHTTPAdapter.WithHandlers(map[string]http.HandlerFunc{
		hla_settings.CHandleIndexPath:          handler.HandleIndexAPI(p.fHTTPLogger),
		hla_settings.CHandleConfigSettingsPath: handler.HandleConfigSettingsAPI(p.fWrapper.GetConfig(), p.fHTTPLogger),
		hla_settings.CHandleConfigConnectsPath: handler.HandleConfigConnectsAPI(pCtx, p.fWrapper, p.fHTTPLogger, networkNode),
		hla_settings.CHandleNetworkOnlinePath:  handler.HandleNetworkOnlineAPI(p.fHTTPLogger, networkNode),
	})
}
