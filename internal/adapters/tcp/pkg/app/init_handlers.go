package app

import (
	"context"

	"github.com/number571/hidden-lake/internal/adapters/tcp/internal/handler"
	hla_http "github.com/number571/hidden-lake/pkg/adapters/http"
	hla_settings "github.com/number571/hidden-lake/pkg/adapters/http/settings"
)

func (p *sApp) initHandlers(pCtx context.Context) {
	networkNode := p.fTCPAdapter.GetConnKeeper().GetNetworkNode()
	p.fHTTPAdapter.
		WithHandlers(
			hla_http.NewHandler(
				hla_settings.CHandleIndexPath,
				handler.HandleIndexAPI(p.fHTTPLogger),
			),
			hla_http.NewHandler(
				hla_settings.CHandleConfigSettingsPath,
				handler.HandleConfigSettingsAPI(p.fWrapper.GetConfig(), p.fHTTPLogger),
			),
			hla_http.NewHandler(
				hla_settings.CHandleConfigConnectsPath,
				handler.HandleConfigConnectsAPI(pCtx, p.fWrapper, p.fHTTPLogger, networkNode),
			),
			hla_http.NewHandler(
				hla_settings.CHandleNetworkOnlinePath,
				handler.HandleNetworkOnlineAPI(p.fHTTPLogger, networkNode),
			),
		)
}
