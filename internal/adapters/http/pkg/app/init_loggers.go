package app

import (
	"github.com/number571/hidden-lake/internal/adapters"
	hla_http_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
)

func (p *sApp) initLoggers() {
	p.fIntAdapter.WithLogger(hla_http_settings.GetAppShortNameFMT(adapters.CAdapterInternalSuffix), p.fAnonLogger)
	p.fExtAdapter.WithLogger(hla_http_settings.GetAppShortNameFMT(), p.fAnonLogger)
}
