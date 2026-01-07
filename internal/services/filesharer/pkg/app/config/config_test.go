package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/hidden-lake/internal/utils/language"
)

const (
	tcLogging    = true
	tcConfigFile = "config_test.txt"
)

const (
	tcConfigTemplate = `settings:
  page_offset: %d
  retry_num: %d
  language: RUS
logging:
  - info
  - erro
address: '%s'
connection: '%s'`
)

const (
	tcAddressIncoming   = "address_incoming"
	tcConnectionService = "connection_service"
	tcMessageSize       = (1 << 20)
	tcPageOffset        = 10
	tcRetryNum          = 2
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
		tcPageOffset,
		tcRetryNum,
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

	if cfg.GetLogging().HasInfo() != tcLogging {
		t.Fatal("logging.info is invalid")
	}

	if cfg.GetLogging().HasErro() != tcLogging {
		t.Fatal("logging.erro is invalid")
	}

	if cfg.GetLogging().HasWarn() == tcLogging {
		t.Fatal("logging.warn is invalid")
	}

	if cfg.GetAddress() != tcAddressIncoming {
		t.Fatal("address.incoming is invalid")
	}

	if cfg.GetConnection() != tcConnectionService {
		t.Fatal("connection.service is invalid")
	}

	if cfg.GetSettings().GetPageOffset() != tcPageOffset {
		t.Fatal("settings.page_offset is invalid")
	}

	if cfg.GetSettings().GetRetryNum() != tcRetryNum {
		t.Fatal("settings.retry_num is invalid")
	}

	if cfg.GetSettings().GetLanguage() != language.CLangRUS {
		t.Fatal("settings language is invalid")
	}
}
