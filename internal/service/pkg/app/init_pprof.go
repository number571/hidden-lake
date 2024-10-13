package app

import (
	"github.com/number571/hidden-lake/internal/modules/pprof"
)

func (p *sApp) initServicePPROF() {
	p.fServicePPROF = pprof.InitPprofService(p.fCfgW.GetConfig().GetAddress().GetPPROF())
}
