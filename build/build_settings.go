// nolint: err113
package build

import (
	_ "embed"
	"errors"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	//go:embed settings.yml
	gSettings []byte
	GSettings SSettings
)

func init() {
	if err := encoding.DeserializeYAML(gSettings, &GSettings); err != nil {
		panic(err)
	}
	if err := GSettings.validate(); err != nil {
		panic(err)
	}
}

type SSettings struct {
	FProtoMask struct {
		FNetwork uint32 `yaml:"network"`
		FService uint32 `yaml:"service"`
	} `yaml:"proto_mask"`
	FQueueProblem struct {
		FMainPoolCap uint64 `yaml:"main_pool_cap"`
		FRandPoolCap uint64 `yaml:"rand_pool_cap"`
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
		p.FQueueProblem.FRandPoolCap == 0:
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

func (p SSettings) GetKeeperPeriod() time.Duration {
	return time.Duration(p.FNetworkManager.FKeeperPeriodMS) * time.Millisecond //nolint:gosec
}

func (p SSettings) GetWriteTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FWriteTimeoutMS) * time.Millisecond //nolint:gosec
}

func (p SSettings) GetReadTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FReadTimeoutMS) * time.Millisecond //nolint:gosec
}

func (p SSettings) GetDialTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FDialTimeoutMS) * time.Millisecond //nolint:gosec
}

func (p SSettings) GetWaitTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FWaitTimeoutMS) * time.Millisecond //nolint:gosec
}
