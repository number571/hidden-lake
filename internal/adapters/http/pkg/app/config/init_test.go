package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/hidden-lake/build"
	hlk_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
)

func TestRebuild(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 99)
	defer func() { _ = os.Remove(configFile) }()

	testConfigDefaultInit(configFile)

	if _, err := InitConfig(configFile, nil, "test_rebuild_config_network"); err == nil {
		t.Fatal("success init config with rebuild for unknown network")
	}

	network := ""
	for k := range build.GetNetworks() {
		network = k
		break
	}

	if _, err := InitConfig(configFile, nil, network); err != nil {
		t.Fatal(err)
	}
}

func TestInit(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 6)
	defer func() { _ = os.Remove(configFile) }()

	testConfigDefaultInit(configFile)

	config1, err := InitConfig(configFile, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	if config1.GetAddress().GetExternal() != tcAddressExternal {
		t.Fatal("got invalid field with exist config (1)")
	}

	_ = os.Remove(configFile)
	if err := os.WriteFile(configFile, []byte("abc"), 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := InitConfig(configFile, nil, ""); err == nil {
		t.Fatal("success init config with invalid config structure (1)")
	}

	_ = os.Remove(configFile)

	config2, err := InitConfig(configFile, config1.(*SConfig), "")
	if err != nil {
		t.Fatal(err)
	}

	if config2.GetAddress().GetExternal() != tcAddressExternal {
		t.Fatal("got invalid field with exist config (2)")
	}

	_ = os.Remove(configFile)

	config3, err := InitConfig(configFile, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	if config3.GetAddress().GetExternal() != hlk_settings.CDefaultExternalAddress {
		t.Fatal("got invalid field with exist config (3)")
	}
}
