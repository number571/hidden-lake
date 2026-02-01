package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/hidden-lake/build"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

func TestRebuild(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 99)
	defer func() { _ = os.Remove(configFile) }()

	testConfigDefaultInit(configFile)

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

	configFile := fmt.Sprintf(tcConfigFileTemplate, 1)
	defer func() { _ = os.Remove(configFile) }()

	testConfigDefaultInit(configFile)

	config1, err := InitConfig(configFile, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(config1.GetApplications()) != 3 {
		t.Fatal("got len(applications) != 3")
	}

	_ = os.Remove(configFile)
	if err := os.WriteFile(configFile, []byte("abc"), 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := InitConfig(configFile, nil, ""); err == nil {
		t.Fatal("success init config with invalid config structure (1)")
	}

	_ = os.Remove(configFile)

	if _, err := InitConfig(configFile, &SConfig{}, ""); err == nil {
		t.Fatal("success init config with invalid config structure (2)")
	}

	_ = os.Remove(configFile)

	config2, err := InitConfig(configFile, config1.(*SConfig), "")
	if err != nil {
		t.Fatal(err)
	}

	if config2.GetApplications()[0] != tgApplications[0] {
		t.Fatal("got invalid field with exist config (2)")
	}

	_ = os.Remove(configFile)

	config3, err := InitConfig(configFile, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	if config3.GetApplications()[0] != hlk_settings.CAppShortName {
		t.Fatal("got invalid field with exist config (3)")
	}
}
