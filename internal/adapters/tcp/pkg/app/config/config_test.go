package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

const (
	tcLogging         = true
	tcMessageSize     = 8192
	tcWorkSize        = 10
	tcNetwork         = "_"
	tcDatabaseEnabled = true
	tcAddressExternal = "external_address"
	tcAddressInternal = "internal_address"
)

var (
	tgEndpoints = []string{
		"endpoint_1",
		"endpoint_2",
	}
	tgConnections = []string{
		"connection_1",
		"connection_2",
	}
)

const (
	tcConfigTemplate = `settings:
  message_size_bytes: %d
  work_size_bits: %d
  network_key: %s
  database_enabled: %t
logging:
  - info
  - erro
address:
  external: %s
  internal: %s
endpoints:
  - %s
  - %s
connections:
  - %s
  - %s
`
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessageSize,
		tcWorkSize,
		tcNetwork,
		tcDatabaseEnabled,
		tcAddressExternal,
		tcAddressInternal,
		tgEndpoints[0],
		tgEndpoints[1],
		tgConnections[0],
		tgConnections[1],
	)
}

func testConfigDefaultInit(configPath string) {
	_ = os.WriteFile(configPath, []byte(testNewConfigString()), 0o600)
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestBuildConfig(t *testing.T) {
	t.Parallel()

	config1File := fmt.Sprintf(tcConfigFileTemplate, 2)
	config2File := fmt.Sprintf(tcConfigFileTemplate, 3)

	testConfigDefaultInit(config1File)
	defer func() { _ = os.Remove(config1File) }()

	cfg, err := LoadConfig(config1File)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := BuildConfig(config2File, cfg.(*SConfig)); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(config2File) }()

	if _, err := BuildConfig(config2File, cfg.(*SConfig)); err == nil {
		t.Fatal("success build already exist config")
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

	cfg1Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "settings", "settings_v2"))
	if err := os.WriteFile(configFile, cfg1Bytes, 0o600); err != nil {
		return err
	}

	cfg2Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "erro", "erro_v2"))
	if err := os.WriteFile(configFile, cfg2Bytes, 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (logging)") // nolint: err113
	}

	return nil
}

func TestComplexConfig(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 0)
	defer func() { _ = os.Remove(configFile) }()

	if err := testIncorrectConfig(configFile); err != nil {
		t.Fatal(err)
	}

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.GetSettings().GetWorkSizeBits() != tcWorkSize {
		t.Fatal("settings work size is invalid")
	}

	if cfg.GetSettings().GetMessageSizeBytes() != tcMessageSize {
		t.Fatal("settings message size is invalid")
	}

	if cfg.GetSettings().GetNetworkKey() != tcNetwork {
		t.Fatal("settings message network_key is invalid")
	}

	if cfg.GetSettings().GetDatabaseEnabled() != tcDatabaseEnabled {
		t.Fatal("settings message database_enabled is invalid")
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

	if cfg.GetSettings().GetNetworkKey() != tcNetwork {
		t.Fatal("network is invalid")
	}

	if cfg.GetAddress().GetExternal() != tcAddressExternal {
		t.Fatal("address_external is invalid")
	}

	if cfg.GetAddress().GetInternal() != tcAddressInternal {
		t.Fatal("address_internal is invalid")
	}

	if len(cfg.GetEndpoints()) != 2 {
		t.Fatal("len connections != 2")
	}
	for i, v := range cfg.GetEndpoints() {
		if v != tgEndpoints[i] {
			t.Fatalf("connection '%d' is invalid", i)
		}
	}
	for i, v := range cfg.GetConnections() {
		if v != tgConnections[i] {
			t.Fatalf("connection '%d' is invalid", i)
		}
	}
}

func TestWrapper(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 1)

	testConfigDefaultInit(configFile)
	defer func() { _ = os.Remove(configFile) }()

	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.GetConnections()) == 0 {
		t.Fatal("list of connections should be is not nil for tests")
	}

	wrapper := NewWrapper(cfg)
	_ = wrapper.GetEditor().UpdateConnections(nil)

	if len(cfg.GetConnections()) != 0 {
		t.Fatal("connections is not nil for current config")
	}

	cfg, err = LoadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.GetConnections()) != 0 {
		t.Fatal("connections is not nil for loaded config")
	}
}
