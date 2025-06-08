package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
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
	tcSendTimeoutMS   = 6_000
	tcRecvTimeoutMS   = 7_000
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
  send_timeout_ms: %d
  recv_timeout_ms: %d
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
		tcSendTimeoutMS,
		tcRecvTimeoutMS,
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
	err := &SConfigError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
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
		t.Error(err)
		return
	}

	if _, err := BuildConfig(config2File, cfg.(*SConfig)); err != nil {
		t.Error(err)
		return
	}
	defer func() { _ = os.Remove(config2File) }()

	if _, err := BuildConfig(config2File, cfg.(*SConfig)); err == nil {
		t.Error("success build already exist config")
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
		t.Error(err)
		return
	}

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	if cfg.GetSettings().GetWorkSizeBits() != tcWorkSize {
		t.Error("settings work size is invalid")
		return
	}

	if cfg.GetSettings().GetMessageSizeBytes() != tcMessageSize {
		t.Error("settings message size is invalid")
		return
	}

	if cfg.GetSettings().GetNetworkKey() != tcNetwork {
		t.Error("settings message network_key is invalid")
		return
	}

	if cfg.GetSettings().GetDatabaseEnabled() != tcDatabaseEnabled {
		t.Error("settings message database_enabled is invalid")
		return
	}

	if cfg.GetSettings().GetSendTimeout() != time.Duration(tcSendTimeoutMS)*time.Millisecond {
		t.Error("settings message send_timeout_ms is invalid")
		return
	}

	if cfg.GetSettings().GetRecvTimeout() != time.Duration(tcRecvTimeoutMS)*time.Millisecond {
		t.Error("settings message recv_timeout_ms is invalid")
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

	if cfg.GetSettings().GetNetworkKey() != tcNetwork {
		t.Error("network is invalid")
		return
	}

	if cfg.GetAddress().GetExternal() != tcAddressExternal {
		t.Error("address_external is invalid")
		return
	}

	if cfg.GetAddress().GetInternal() != tcAddressInternal {
		t.Error("address_internal is invalid")
		return
	}

	if len(cfg.GetEndpoints()) != 2 {
		t.Error("len connections != 2")
		return
	}
	for i, v := range cfg.GetEndpoints() {
		if v != tgEndpoints[i] {
			t.Errorf("connection '%d' is invalid", i)
			return
		}
	}
	for i, v := range cfg.GetConnections() {
		if v != tgConnections[i] {
			t.Errorf("connection '%d' is invalid", i)
			return
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
		t.Error(err)
		return
	}

	if len(cfg.GetConnections()) == 0 {
		t.Error("list of connections should be is not nil for tests")
		return
	}

	wrapper := NewWrapper(cfg)
	_ = wrapper.GetEditor().UpdateConnections(nil)

	if len(cfg.GetConnections()) != 0 {
		t.Error("connections is not nil for current config")
		return
	}

	cfg, err = LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	if len(cfg.GetConnections()) != 0 {
		t.Error("connections is not nil for loaded config")
		return
	}
}
