package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/internal/adapters/http/internal/handler"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
)

func (p *sApp) initHandlers(pCtx context.Context) {
	p.fHTTPIntAdapter.WithHandlers(map[string]http.HandlerFunc{
		hla_settings.CHandleIndexPath:          handler.HandleIndexAPI(p.fHTTPLogger),
		hla_settings.CHandleConfigSettingsPath: handler.HandleConfigSettingsAPI(p.fWrapper.GetConfig(), p.fHTTPLogger),
		hla_settings.CHandleConfigConnectsPath: handler.HandleConfigConnectsAPI(pCtx, p.fWrapper, p.fHTTPLogger),
		hla_settings.CHandleNetworkOnlinePath:  handler.HandleNetworkOnlineAPI(p.fHTTPLogger, p.fHTTPExtAdapter),
	})
}
