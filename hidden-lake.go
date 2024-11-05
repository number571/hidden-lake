package hiddenlake

import (
	_ "embed"
	"fmt"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	CVersion        = "v1.7.4"
	CDefaultNetwork = "__default_network__"
)

var (
	//go:embed networks.yml
	gNetworks []byte
	GNetworks map[string]SNetwork
)

func init() {
	networksYAML := &SNetworksYAML{}
	if err := encoding.DeserializeYAML(gNetworks, networksYAML); err != nil {
		panic(err)
	}
	if _, ok := networksYAML.FNetworks[CDefaultNetwork]; ok {
		panic(fmt.Sprintf("network '%s' already exist", CDefaultNetwork))
	}
	GNetworks = networksYAML.FNetworks
	GNetworks[CDefaultNetwork] = networksYAML.FSettings
}

type SNetworksYAML struct {
	FSettings SNetwork            `yaml:"settings"`
	FNetworks map[string]SNetwork `yaml:"networks"`
}

type SNetwork struct {
	FMessageSizeBytes uint64        `yaml:"message_size_bytes"`
	FFetchTimeoutMS   uint64        `yaml:"fetch_timeout_ms"`
	FQueuePeriodMS    uint64        `yaml:"queue_period_ms"`
	FWorkSizeBits     uint64        `yaml:"work_size_bits"`
	FConnections      []SConnection `yaml:"connections"`
}

type SConnection struct {
	FHost string `yaml:"host"`
	FPort uint16 `yaml:"port"`
}
