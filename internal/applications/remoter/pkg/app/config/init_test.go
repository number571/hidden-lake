package config

import (
	"fmt"
	"os"
	"testing"

	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

func TestInit(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 1)
	defer func() { _ = os.Remove(configFile) }()

	testConfigDefaultInit(configFile)

	config1, err := InitConfig(configFile, nil)
	if err != nil {
		t.Fatal(err)
	}

	if config1.GetAddress().GetExternal() != tcAddress1 {
		t.Fatal("got invalid field with exist config (1)")
	}

	_ = os.Remove(configFile)
	if err := os.WriteFile(configFile, []byte("abc"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := InitConfig(configFile, nil); err == nil {
		t.Fatal("success init config with invalid config structure (1)")
	}

	_ = os.Remove(configFile)

	if _, err := InitConfig(configFile, &SConfig{}); err == nil {
		t.Fatal("success init config with invalid config structure (2)")
	}

	_ = os.Remove(configFile)

	config3, err := InitConfig(configFile, nil)
	if err != nil {
		t.Fatal(err)
	}

	if config3.GetAddress().GetExternal() != hlr_settings.CDefaultExternalAddress {
		t.Fatal("got invalid field with exist config (3)")
	}
}
