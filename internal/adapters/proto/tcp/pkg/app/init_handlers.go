package app

import (
	"context"

	hla_settings "github.com/number571/hidden-lake/internal/adapters/pkg/settings"
	"github.com/number571/hidden-lake/internal/adapters/proto/tcp/internal/handler"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/proto/tcp/pkg/settings"
	hla_http "github.com/number571/hidden-lake/pkg/adapters/http"
)

func (p *sApp) initHandlers(pCtx context.Context) {
	p.fHTTPAdapter.
		WithLogger(hla_tcp_settings.GServiceName, p.fHTTPLogger).
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
				handler.HandleConfigConnectsAPI(pCtx, p.fWrapper, p.fHTTPLogger, p.fTCPAdapter),
			),
			hla_http.NewHandler(
				hla_settings.CHandleNetworkOnlinePath,
				handler.HandleNetworkOnlineAPI(p.fHTTPLogger, p.fTCPAdapter),
			),
		)
}
