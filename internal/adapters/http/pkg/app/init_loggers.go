package app

import hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"

func (p *sApp) initLoggers() {
	p.fHTTPIntAdapter.WithLogger(hla_tcp_settings.GetServiceName(), p.fAnonLogger)
}
