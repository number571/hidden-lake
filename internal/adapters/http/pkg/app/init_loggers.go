package app

import hla_http_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"

func (p *sApp) initLoggers() {
	p.fIntAdapter.WithLogger(hla_http_settings.GetAppShortNameFMT()+"(INT)", p.fAnonLogger)
	p.fExtAdapter.WithLogger(hla_http_settings.GetAppShortNameFMT(), p.fAnonLogger)
}
