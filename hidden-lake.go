// nolint: err113
package hiddenlake

import (
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	CVersion        = "v1.7.5~"
	CDefaultNetwork = "__default_network__"
)

var (
	//go:embed build/networks.yml
	gNetworks []byte
	GNetworks map[string]SNetwork

	//go:embed build/settings.yml
	gSettings []byte
	GSettings SSettings
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

func init() {
	if err := encoding.DeserializeYAML(gSettings, &GSettings); err != nil {
		panic(err)
	}
	if err := GSettings.validate(); err != nil {
		panic(err)
	}
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

type SSettings struct {
	ProtoMask struct {
		Network uint32 `yaml:"network"`
		Service uint32 `yaml:"service"`
	} `yaml:"proto_mask"`
	QueueCapacity struct {
		Main uint64 `yaml:"main"`
		Rand uint64 `yaml:"rand"`
	} `yaml:"queue_capacity"`
	NetworkManager struct {
		ConnectsLimiter uint64 `yaml:"connects_limiter"`
		CacheHashesCap  uint64 `yaml:"cache_hashes_cap"`
		KeeperPeriodMS  uint64 `yaml:"keeper_period_ms"`
	} `yaml:"network_manager"`
	NetworkConnection struct {
		WriteTimeoutMS uint64 `yaml:"write_timeout_ms"`
		ReadTimeoutMS  uint64 `yaml:"read_timeout_ms"`
		DialTimeoutMS  uint64 `yaml:"dial_timeout_ms"`
		WaitTimeoutMS  uint64 `yaml:"wait_timeout_ms"`
	} `yaml:"network_connection"`
}

func (p SSettings) validate() error {
	switch {
	case
		p.QueueCapacity.Main == 0,
		p.QueueCapacity.Rand == 0:
		return errors.New("queue_capacity is invalid")
	case
		p.NetworkManager.ConnectsLimiter == 0,
		p.NetworkManager.CacheHashesCap == 0,
		p.NetworkManager.KeeperPeriodMS == 0:
		return errors.New("network_manager is invalid")
	case
		p.NetworkConnection.WriteTimeoutMS == 0,
		p.NetworkConnection.ReadTimeoutMS == 0,
		p.NetworkConnection.DialTimeoutMS == 0,
		p.NetworkConnection.WaitTimeoutMS == 0:
		return errors.New("network_connection is invalid")
	}
	return nil
}

func (p SSettings) GetKeeperPeriod() time.Duration {
	return time.Duration(p.NetworkManager.KeeperPeriodMS) * time.Millisecond
}

func (p SSettings) GetWriteTimeout() time.Duration {
	return time.Duration(p.NetworkConnection.WriteTimeoutMS) * time.Millisecond
}

func (p SSettings) GetReadTimeout() time.Duration {
	return time.Duration(p.NetworkConnection.ReadTimeoutMS) * time.Millisecond
}

func (p SSettings) GetDialTimeout() time.Duration {
	return time.Duration(p.NetworkConnection.DialTimeoutMS) * time.Millisecond
}

func (p SSettings) GetWaitTimeout() time.Duration {
	return time.Duration(p.NetworkConnection.WaitTimeoutMS) * time.Millisecond
}
