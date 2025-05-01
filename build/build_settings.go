// nolint: err113
package build

import (
	_ "embed"
	"errors"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	//go:embed settings.yml
	gSettingsVal []byte
	gSettingsMtx sync.RWMutex
	gSettings    SSettings
)

func init() {
	settingsYAML := &SSettings{}
	if err := encoding.DeserializeYAML(gSettingsVal, settingsYAML); err != nil {
		panic(err)
	}
	if err := SetSettings(*settingsYAML); err != nil {
		panic(err)
	}
}

func GetSettings() SSettings {
	gSettingsMtx.RLock()
	defer gSettingsMtx.RUnlock()

	return gSettings
}

func SetSettings(settings SSettings) error {
	if err := settings.validate(); err != nil {
		return err
	}
	gSettingsMtx.Lock()
	gSettings = settings
	gSettingsMtx.Unlock()
	return nil
}

type SSettings struct {
	FProtoMask struct {
		FNetwork uint32 `yaml:"network"`
		FService uint32 `yaml:"service"`
	} `yaml:"proto_mask"`
	FQueueBasedProblem struct {
		FMainPoolCap  uint64 `yaml:"main_pool_cap"`
		FRandPoolCap  uint64 `yaml:"rand_pool_cap"`
		FPowParallels uint64 `yaml:"pow_parallels"`
		FQBPConsumers uint64 `yaml:"num_consumers"`
	} `yaml:"queue_based_problem"`
	FNetworkManager struct {
		FConnectsLimiter uint64 `yaml:"connects_limiter"`
		FCacheHashesCap  uint64 `yaml:"cache_hashes_cap"`
		FKeeperPeriodMS  uint64 `yaml:"keeper_period_ms"`
	} `yaml:"network_manager"`
	FNetworkConnection struct {
		FSendTimeoutMS uint64 `yaml:"send_timeout_ms"`
		FRecvTimeoutMS uint64 `yaml:"recv_timeout_ms"`
		FDialTimeoutMS uint64 `yaml:"dial_timeout_ms"`
		FWaitTimeoutMS uint64 `yaml:"wait_timeout_ms"`
	} `yaml:"network_connection"`
}

func (p SSettings) validate() error {
	switch {
	case
		p.FQueueBasedProblem.FMainPoolCap == 0,
		p.FQueueBasedProblem.FRandPoolCap == 0,
		p.FQueueBasedProblem.FQBPConsumers == 0,
		p.FQueueBasedProblem.FPowParallels == 0:
		return errors.New("queue_problem is invalid")
	case
		p.FNetworkManager.FConnectsLimiter == 0,
		p.FNetworkManager.FCacheHashesCap == 0,
		p.FNetworkManager.FKeeperPeriodMS == 0:
		return errors.New("network_manager is invalid")
	case
		p.FNetworkConnection.FSendTimeoutMS == 0,
		p.FNetworkConnection.FRecvTimeoutMS == 0,
		p.FNetworkConnection.FDialTimeoutMS == 0,
		p.FNetworkConnection.FWaitTimeoutMS == 0:
		return errors.New("network_connection is invalid")
	}
	return nil
}

func (p SSettings) GetKeeperPeriod() time.Duration {
	return time.Duration(p.FNetworkManager.FKeeperPeriodMS) * time.Millisecond //nolint:gosec
}

func (p SSettings) GetSendTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FSendTimeoutMS) * time.Millisecond //nolint:gosec
}

func (p SSettings) GetRecvTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FRecvTimeoutMS) * time.Millisecond //nolint:gosec
}

func (p SSettings) GetDialTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FDialTimeoutMS) * time.Millisecond //nolint:gosec
}

func (p SSettings) GetWaitTimeout() time.Duration {
	return time.Duration(p.FNetworkConnection.FWaitTimeoutMS) * time.Millisecond //nolint:gosec
}
