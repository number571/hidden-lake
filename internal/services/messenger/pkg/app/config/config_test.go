package config

import (
	"fmt"
	"os"
	"testing"
)

const (
	tcLogging    = true
	tcConfigFile = "config_test.txt"
)

const (
	tcConfigTemplate = `settings:
  messages_capacity: %d
logging:
  - info
  - erro
address:
  internal: '%s'
  external: '%s'
connection: '%s'`
)

const (
	tcAddressInterface  = "address_interface"
	tcAddressIncoming   = "address_incoming"
	tcConnectionService = "connection_service"
	tcMessagesCapacity  = 1000
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SConfigError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessagesCapacity,
		tcAddressInterface,
		tcAddressIncoming,
		tcConnectionService,
	)
}

func testConfigDefaultInit(configPath string) {
	_ = os.WriteFile(configPath, []byte(testNewConfigString()), 0o600)
}

func TestConfig(t *testing.T) {
	t.Parallel()

	testConfigDefaultInit(tcConfigFile)
	defer func() { _ = os.Remove(tcConfigFile) }()

	cfg, err := LoadConfig(tcConfigFile)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.GetSettings().GetMessagesCapacity() != tcMessagesCapacity {
		t.Fatal("settings message capacity size is invalid")
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

	if cfg.GetAddress().GetInternal() != tcAddressInterface {
		t.Fatal("address.interface is invalid")
	}

	if cfg.GetAddress().GetExternal() != tcAddressIncoming {
		t.Fatal("address.incoming is invalid")
	}

	if cfg.GetConnection() != tcConnectionService {
		t.Fatal("connection.service is invalid")
	}

}
