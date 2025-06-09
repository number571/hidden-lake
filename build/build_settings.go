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
		FCacheHashesCap      uint64 `yaml:"cache_hashes_cap"`
		FHttpReadTimeoutMS   uint64 `yaml:"http_read_timeout_ms"`
		FHttpWriteTimeoutMS  uint64 `yaml:"http_write_timeout_ms"`
		FHttpHandleTimeoutMS uint64 `yaml:"http_handle_timeout_ms"`
	} `yaml:"network_manager"`
	FQueueBasedProblem struct {
		FPoolCap [2]uint64 `yaml:"pool_cap"`
	} `yaml:"queue_based_problem"`
}

func (p SSettings) validate() error {
	switch {
	case
		p.FNetworkManager.FCacheHashesCap == 0,
		p.FNetworkManager.FHttpReadTimeoutMS == 0,
		p.FNetworkManager.FHttpWriteTimeoutMS == 0,
		p.FNetworkManager.FHttpHandleTimeoutMS == 0:
		return errors.New("network_manager is invalid")
	case
		p.FQueueBasedProblem.FPoolCap[0] == 0,
		p.FQueueBasedProblem.FPoolCap[1] == 0:
		return errors.New("queue_based_problem is invalid")
	}
	return nil
}

func (p SSettings) GetHttpReadTimeout() time.Duration {
	return time.Duration(p.FNetworkManager.FHttpReadTimeoutMS) * time.Millisecond // nolint: gosec
}

func (p SSettings) GetHttpWriteTimeout() time.Duration {
	return time.Duration(p.FNetworkManager.FHttpWriteTimeoutMS) * time.Millisecond // nolint: gosec
}

func (p SSettings) GetHttpHandleTimeout() time.Duration {
	return time.Duration(p.FNetworkManager.FHttpHandleTimeoutMS) * time.Millisecond // nolint: gosec
}
