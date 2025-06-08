// nolint: err113
package build

import (
	_ "embed"
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	//go:embed settings.yml
	GSettingsVal []byte
	gSettingsMtx sync.RWMutex
	gSettings    SSettings
)

func init() {
	settingsYAML := &SSettings{}
	if err := encoding.DeserializeYAML(GSettingsVal, settingsYAML); err != nil {
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
	FNetworkManager struct {
		FCacheHashesCap uint64 `yaml:"cache_hashes_cap"`
	} `yaml:"network_manager"`
	FQueueBasedProblem struct {
		FMainPoolCap uint64 `yaml:"main_pool_cap"`
		FRandPoolCap uint64 `yaml:"rand_pool_cap"`
	} `yaml:"queue_based_problem"`
}

func (p SSettings) validate() error {
	switch {
	case
		p.FNetworkManager.FCacheHashesCap == 0:
		return errors.New("network_manager is invalid")
	case
		p.FQueueBasedProblem.FMainPoolCap == 0,
		p.FQueueBasedProblem.FRandPoolCap == 0:
		return errors.New("queue_based_problem is invalid")
	}
	return nil
}
