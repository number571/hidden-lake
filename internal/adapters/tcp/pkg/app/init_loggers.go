package app

import (
	"github.com/number571/hidden-lake/internal/adapters"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
)

func (p *sApp) initLoggers() {
	p.fIntAdapter.WithLogger(hla_tcp_settings.GetAppShortNameFMT(adapters.CAdapterInternalSuffix), p.fAnonLogger)
	p.fExtAdapter.WithLogger(hla_tcp_settings.GetAppShortNameFMT(), p.fAnonLogger)
}
