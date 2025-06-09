package config

import (
	"os"
	"testing"
)

const (
	tcExecTimeout = 5000
	tcConfigFile  = "config_test.txt"
	tcLogging     = true
	tcAddress1    = "test_address1"
	tcAddress2    = "test_address2"
	tcPassword    = "test_password"
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
		FSettings: &SConfigSettings{
			FPassword: tcPassword,
		},
		FLogging: []string{"info", "erro"},
		FAddress: &SAddress{
			FExternal: tcAddress1,
		},
	})
}

func TestConfig(t *testing.T) {
	t.Parallel()

	testConfigDefaultInit(tcConfigFile)
	defer func() { _ = os.Remove(tcConfigFile) }()

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Error(err)
		return
	}

	if cfg.GetLogging().HasInfo() != tcLogging {
		t.Error("logging.info is invalid")
		return
	}

	if cfg.GetLogging().HasErro() != tcLogging {
		t.Error("logging.erro is invalid")
		return
	}

	if cfg.GetLogging().HasWarn() == tcLogging {
		t.Error("logging.warn is invalid")
		return
	}

	if cfg.GetAddress().GetExternal() != tcAddress1 {
		t.Error("address incoming is invalid")
		return
	}

	if cfg.GetSettings().GetPassword() != tcPassword {
		t.Error("settings.password is invalid")
		return
	}
}
