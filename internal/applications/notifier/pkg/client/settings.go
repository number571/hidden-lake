package client

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FDiffBits uint64
	FParallel uint64
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
	}
	return (&sSettings{
		FDiffBits: pSett.FDiffBits,
		FParallel: pSett.FParallel,
	}).useDefault()
}

func (p *sSettings) useDefault() *sSettings {
	if p.FParallel == 0 {
		p.FParallel = 1
	}
	return p
}

func (p *sSettings) GetDiffBits() uint64 {
	return p.FDiffBits
}

func (p *sSettings) GetParallel() uint64 {
	return p.FParallel
}
