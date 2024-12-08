package http

import (
	"github.com/number571/hidden-lake/build"
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FAddress          string
	FProducePath      string
	FMessageSizeBytes uint64
	FWorkSizeBits     uint64
	FNetworkKey       string
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{}
	}
	return (&sSettings{
		FAddress:          pSett.FAddress,
		FProducePath:      pSett.FProducePath,
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
	if p.FProducePath == "" {
		p.FProducePath = "/"
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FAddress
}

func (p *sSettings) GetProducePath() string {
	return p.FProducePath
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
