package config

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

const (
	tcConfigFile = "config_test.txt"
)

var (
	tgServices = []string{"service_1", "service_2", "service_3"}
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SConfigError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func testConfigDefaultInit(configPath string) {
	_, _ = BuildConfig(configPath, &SConfig{
		FLogging:  []string{"info", "erro"},
		FServices: tgServices,
	})
}

func TestConfig(t *testing.T) {
	t.Parallel()

	testConfigDefaultInit(tcConfigFile)
	defer os.Remove(tcConfigFile)

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
		return
	}

	if cfg.GetLogging().HasInfo() != true {
		t.Error("invalid logging info")
		return
	}

	if cfg.GetLogging().HasWarn() != false {
		t.Error("invalid logging warn")
		return
	}

	if cfg.GetLogging().HasErro() != true {
		t.Error("invalid logging erro")
		return
	}

	services := cfg.GetServices()
	if len(services) != 3 {
		t.Error("settings value is invalid")
		return
	}

	for i := range services {
		if services[i] != tgServices[i] {
			t.Error("got invalid service")
			return
		}
	}
}

func TestComplexConfig(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 0)
	defer os.Remove(configFile)

	if err := testIncorrectConfig(configFile); err != nil {
		t.Error(err)
		return
	}
}

func testIncorrectConfig(configFile string) error {
	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config on non exist file") // nolint: err113
	}

	if err := os.WriteFile(configFile, []byte("abc"), 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid structure") // nolint: err113
	}

	return nil
}
