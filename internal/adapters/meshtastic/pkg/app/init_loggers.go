package app

import (
	"github.com/number571/hidden-lake/internal/adapters"
	hla_meshtastic_settings "github.com/number571/hidden-lake/internal/adapters/meshtastic/pkg/settings"
)

func (p *sApp) initLoggers() {
	p.fIntAdapter.WithLogger(hla_meshtastic_settings.GetAppShortNameFMT(adapters.CAdapterInternalSuffix), p.fAnonLogger)
	p.fExtAdapter.WithLogger(hla_meshtastic_settings.GetAppShortNameFMT(), p.fAnonLogger)
}
