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

func NewSettingsByNetworkKey(pNetworkKey string) ISettings {
	network, ok := build.GetNetwork(pNetworkKey)
	if !ok {
		panic("network key undefined")
	}
	return NewSettings(&SSettings{
		FNetworkKey:       pNetworkKey,
		FWorkSizeBits:     network.FWorkSizeBits,
		FMessageSizeBytes: network.FMessageSizeBytes,
	})
}

func (p *sSettings) useDefault() *sSettings {
	defaultNetwork, _ := build.GetNetwork(build.CDefaultNetwork)
	if p.FMessageSizeBytes == 0 {
		p.FMessageSizeBytes = defaultNetwork.FMessageSizeBytes
	}
	if p.FWorkSizeBits == 0 {
		p.FWorkSizeBits = defaultNetwork.FWorkSizeBits
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
