package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/hidden-lake/build"
	hls_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
)

func TestRebuild(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 99)
	defer os.Remove(configFile)

	testConfigDefaultInit(configFile)

	if _, err := InitConfig(configFile, nil, "test_rebuild_config_network"); err == nil {
		t.Error("success init config with rebuild for unknown network")
		return
	}

	network := ""
	for k := range build.GNetworks {
		network = k
		break
	}

	if _, err := InitConfig(configFile, nil, network); err != nil {
		t.Error(err)
		return
	}
}

func TestInit(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 6)
	defer os.Remove(configFile)

	testConfigDefaultInit(configFile)

	config1, err := InitConfig(configFile, nil, "")
	if err != nil {
		t.Error(err)
		return
	}

	if config1.GetAddress().GetExternal() != tcAddressExternal {
		t.Error("got invalid field with exist config (1)")
		return
	}

	os.Remove(configFile)
	if err := os.WriteFile(configFile, []byte("abc"), 0o600); err != nil {
		t.Error(err)
		return
	}

	if _, err := InitConfig(configFile, nil, ""); err == nil {
		t.Error("success init config with invalid config structure (1)")
		return
	}

	os.Remove(configFile)

	if _, err := InitConfig(configFile, &SConfig{}, ""); err == nil {
		t.Error("success init config with invalid config structure (2)")
		return
	}

	os.Remove(configFile)

	config2, err := InitConfig(configFile, config1.(*SConfig), "")
	if err != nil {
		t.Error(err)
		return
	}

	if config2.GetAddress().GetExternal() != tcAddressExternal {
		t.Error("got invalid field with exist config (2)")
		return
	}

	os.Remove(configFile)

	config3, err := InitConfig(configFile, nil, "")
	if err != nil {
		t.Error(err)
		return
	}

	if config3.GetAddress().GetExternal() != hls_settings.CDefaultExternalAddress {
		t.Error("got invalid field with exist config (3)")
		return
	}
}
