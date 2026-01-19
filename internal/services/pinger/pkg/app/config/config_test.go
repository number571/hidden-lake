package config

import (
	"os"
	"testing"
)

const (
	tcConfigFile = "config_test.txt"
	tcLogging    = true
	tcAddress1   = "test_address1"
	tcAddress2   = "test_address2"
	tcAddress3   = "test_address3"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func testConfigDefaultInit(configPath string) {
	_, _ = BuildConfig(configPath, &SConfig{
		FSettings: &SConfigSettings{},
		FLogging:  []string{"info", "erro"},
		FAddress: &SAddress{
			FInternal: tcAddress2,
			FExternal: tcAddress1,
		},
		FConnection: tcAddress3,
	})
}

func TestConfig(t *testing.T) {
	t.Parallel()

	testConfigDefaultInit(tcConfigFile)
	defer func() { _ = os.Remove(tcConfigFile) }()

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.GetLogging().HasInfo() != tcLogging {
		t.Fatal("logging.info is invalid")
	}

	if cfg.GetLogging().HasErro() != tcLogging {
		t.Fatal("logging.erro is invalid")
	}

	if cfg.GetLogging().HasWarn() == tcLogging {
		t.Fatal("logging.warn is invalid")
	}

	if cfg.GetAddress().GetExternal() != tcAddress1 {
		t.Fatal("address incoming is invalid")
	}
}
