package app

import "github.com/number571/hidden-lake/internal/utils/pprof"

func (p *sApp) initServicePPROF() {
	p.fServicePPROF = pprof.InitPprofService(p.fWrapper.GetConfig().GetAddress().GetPPROF())
}
