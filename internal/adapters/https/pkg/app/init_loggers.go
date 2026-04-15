package app

import hla_https_settings "github.com/number571/hidden-lake/internal/adapters/https/pkg/settings"

func (p *sApp) initLoggers() {
	p.fIntAdapter.WithLogger(hla_https_settings.GetAppShortNameFMT()+"(INT)", p.fAnonLogger)
	p.fExtAdapter.WithLogger(hla_https_settings.GetAppShortNameFMT(), p.fAnonLogger)
}
