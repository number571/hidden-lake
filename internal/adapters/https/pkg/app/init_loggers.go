package app

import (
	"github.com/number571/hidden-lake/internal/adapters"
	hla_https_settings "github.com/number571/hidden-lake/internal/adapters/https/pkg/settings"
)

func (p *sApp) initLoggers() {
	p.fIntAdapter.WithLogger(hla_https_settings.GetAppShortNameFMT(adapters.CAdapterInternalSuffix), p.fAnonLogger)
	p.fExtAdapter.WithLogger(hla_https_settings.GetAppShortNameFMT(), p.fAnonLogger)
}
