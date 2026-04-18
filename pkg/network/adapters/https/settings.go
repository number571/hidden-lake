package https

import (
	"time"

	"github.com/number571/hidden-lake/pkg/network/adapters"
)

const (
	CDefaultRateLimit     = 5
	CDefaultCapacityLimit = 10
	CDefaultChannelSize   = 64
	CDefaultReadTimeout   = 5 * time.Second
	CDefaultHandleTimeout = 5 * time.Second
)

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FAdapterSettings adapters.ISettings
	FServeSettings   *SServeSettings
}

type SServeSettings struct {
	FAddress         string
	FAuthMapper      map[string]string
	FRateLimitParams [2]float64
	FDataBrokerParam uint64
	FReadTimeout     time.Duration
	FHandleTimeout   time.Duration
}

func NewSettings(pSett *SSettings) ISettings {
	if pSett == nil {
		pSett = &SSettings{
			FAdapterSettings: adapters.NewSettings(nil),
		}
	}
	return (&sSettings{
		FAdapterSettings: pSett.FAdapterSettings,
		FServeSettings:   pSett.FServeSettings,
	}).initDefault()
}

func (p *sSettings) initDefault() *sSettings {
	if p.FServeSettings == nil {
		p.FServeSettings = &SServeSettings{}
	}
	if p.FServeSettings.FRateLimitParams[0] <= 0 {
		p.FServeSettings.FRateLimitParams[0] = CDefaultRateLimit
	}
	if p.FServeSettings.FRateLimitParams[1] <= 0 {
		p.FServeSettings.FRateLimitParams[1] = CDefaultCapacityLimit
	}
	if p.FServeSettings.FDataBrokerParam == 0 {
		p.FServeSettings.FDataBrokerParam = CDefaultChannelSize
	}
	if p.FServeSettings.FReadTimeout == 0 {
		p.FServeSettings.FReadTimeout = CDefaultReadTimeout
	}
	if p.FServeSettings.FHandleTimeout == 0 {
		p.FServeSettings.FHandleTimeout = CDefaultHandleTimeout
	}
	return p
}

func (p *sSettings) GetAddress() string {
	return p.FServeSettings.FAddress
}

func (p *sSettings) GetAuthMapper() map[string]string {
	return p.FServeSettings.FAuthMapper
}

func (p *sSettings) GetRateLimitParams() [2]float64 {
	return p.FServeSettings.FRateLimitParams
}

func (p *sSettings) GetDataBrokerParam() uint64 {
	return p.FServeSettings.FDataBrokerParam
}

func (p *sSettings) GetAdapterSettings() adapters.ISettings {
	return p.FAdapterSettings
}

func (p *sSettings) GetReadTimeout() time.Duration {
	return p.FServeSettings.FReadTimeout
}

func (p *sSettings) GetHandleTimeout() time.Duration {
	return p.FServeSettings.FHandleTimeout
}
