package app

import (
	"github.com/number571/hidden-lake/internal/pprof"
)

func (p *sApp) initServicePPROF() {
	p.fServicePPROF = pprof.InitPprofService(p.fConfig.GetAddress().GetPPROF())
}
