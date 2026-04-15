package app

import hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"

func (p *sApp) initLoggers() {
	p.fIntAdapter.WithLogger(hla_tcp_settings.GetAppShortNameFMT(), p.fAnonLogger)
	p.fExtAdapter.WithLogger(hla_tcp_settings.GetAppShortNameFMT(), p.fAnonLogger)
}
