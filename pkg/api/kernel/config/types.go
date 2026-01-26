package config

import "github.com/number571/hidden-lake/internal/kernel/pkg/app/config"

type IConfigSettings interface {
	config.IConfigSettings
	GetPayloadSizeBytes() uint64
}

type SConfigSettings struct {
	config.SConfigSettings
	FPayloadSizeBytes uint64 `json:"payload_size_bytes"`
}

func (p *SConfigSettings) GetPayloadSizeBytes() uint64 {
	return p.FPayloadSizeBytes
}
