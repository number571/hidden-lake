package client

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FWorkSizeBits uint64
	FPowParallel  uint64
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
	}
	return (&sSettings{
		FWorkSizeBits: pSett.FWorkSizeBits,
		FPowParallel:  pSett.FPowParallel,
	}).useDefault()
}

func (p *sSettings) useDefault() *sSettings {
	if p.FPowParallel == 0 {
		p.FPowParallel = 1
	}
	return p
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetPowParallel() uint64 {
	return p.FPowParallel
}
