package adapters

import (
	"github.com/number571/hidden-lake/build"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FMessageSizeBytes uint64
	FWorkSizeBits     uint64
	FNetworkKey       string
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
	}
	return (&sSettings{
		FMessageSizeBytes: pSett.FMessageSizeBytes,
		FWorkSizeBits:     pSett.FWorkSizeBits,
		FNetworkKey:       pSett.FNetworkKey,
	}).useDefault()
}

func (p *sSettings) useDefault() *sSettings {
	defaultNetwork := build.GNetworks[build.CDefaultNetwork]
	if p.FMessageSizeBytes == 0 {
		p.FMessageSizeBytes = defaultNetwork.FMessageSizeBytes
	}
	return p
}

func (p *sSettings) GetNetworkKey() string {
	return p.FNetworkKey
}

func (p *sSettings) GetWorkSizeBits() uint64 {
	return p.FWorkSizeBits
}

func (p *sSettings) GetMessageSizeBytes() uint64 {
	return p.FMessageSizeBytes
}
