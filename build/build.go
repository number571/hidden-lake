// nolint: err113
package build

import (
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	CDefaultNetwork = "__default_network__"
)

var (
	//go:embed networks.yml
	gNetworks []byte
	GNetworks map[string]SNetwork

	//go:embed settings.yml
	gSettings []byte
	GSettings SSettings

	//go:embed version.yml
	gVersion []byte
	GVersion string
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

func init() {
	var versionYAML struct {
		FVersion string `yaml:"version"`
	}
	if err := encoding.DeserializeYAML(gVersion, &versionYAML); err != nil {
		panic(err)
	}
	GVersion = versionYAML.FVersion
}

type SNetworksYAML struct {
	FSettings SNetwork            `yaml:"settings"`
	FNetworks map[string]SNetwork `yaml:"networks"`
}

type SNetwork struct {
	FMessageSizeBytes uint64   `yaml:"message_size_bytes"`
	FFetchTimeoutMS   uint64   `yaml:"fetch_timeout_ms"`
	FQueuePeriodMS    uint64   `yaml:"queue_period_ms"`
	FWorkSizeBits     uint64   `yaml:"work_size_bits"`
	FConnections      []string `yaml:"connections"`
}

type SSettings struct {
	FProtoMask struct {
		FNetwork uint32 `yaml:"network"`
		FService uint32 `yaml:"service"`
	} `yaml:"proto_mask"`
	FQueueProblem struct {
		FMainPoolCap  uint64 `yaml:"main_pool_cap"`
		FRandPoolCap  uint64 `yaml:"rand_pool_cap"`
		FConsumersCap uint64 `yaml:"consumers_cap"`
	} `yaml:"queue_problem"`
	FNetworkManager struct {
		FConnectsLimiter uint64 `yaml:"connects_limiter"`
		FCacheHashesCap  uint64 `yaml:"cache_hashes_cap"`
		FKeeperPeriodMS  uint64 `yaml:"keeper_period_ms"`
	} `yaml:"network_manager"`
	FNetworkConnection struct {
		FWriteTimeoutMS uint64 `yaml:"write_timeout_ms"`
		FReadTimeoutMS  uint64 `yaml:"read_timeout_ms"`
		FDialTimeoutMS  uint64 `yaml:"dial_timeout_ms"`
		FWaitTimeoutMS  uint64 `yaml:"wait_timeout_ms"`
	} `yaml:"network_connection"`
}

func (p SSettings) validate() error {
	switch {
	case
		p.FQueueProblem.FMainPoolCap == 0,
		p.FQueueProblem.FRandPoolCap == 0,
		p.FQueueProblem.FConsumersCap == 0:
		return errors.New("queue_problem is invalid")
	case
		p.FNetworkManager.FConnectsLimiter == 0,
		p.FNetworkManager.FCacheHashesCap == 0,
		p.FNetworkManager.FKeeperPeriodMS == 0:
		return errors.New("network_manager is invalid")
	case
		p.FNetworkConnection.FWriteTimeoutMS == 0,
		p.FNetworkConnection.FReadTimeoutMS == 0,
		p.FNetworkConnection.FDialTimeoutMS == 0,
		p.FNetworkConnection.FWaitTimeoutMS == 0:
		return errors.New("network_connection is invalid")
	}
	return nil
}

func (p SNetwork) GetFetchTimeout() time.Duration {
	return time.Duration(p.FFetchTimeoutMS) * time.Millisecond
}

func (p SNetwork) GetQueuePeriod() time.Duration {
	return time.Duration(p.FQueuePeriodMS) * time.Millisecond
}

func (p SSettings) GetKeeperPeriod() time.Duration {
	return time.Duration(p.FNetworkManager.FKeeperPeriodMS) * time.Millisecond
}

func (p SSettings) GetWriteTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FWriteTimeoutMS) * time.Millisecond
}

func (p SSettings) GetReadTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FReadTimeoutMS) * time.Millisecond
}

func (p SSettings) GetDialTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FDialTimeoutMS) * time.Millisecond
}

func (p SSettings) GetWaitTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FWaitTimeoutMS) * time.Millisecond
}
