// nolint: err113
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

const (
	tcLogging         = true
	tcNetwork         = "test_network_key"
	tcDownloader      = "test_downloader"
	tcUploader        = "test_uploader"
	tcAddressExternal = "test_address_external"
	tcAddressInternal = "test_address_internal"
	tcPubKeyAlias1    = "test_alias1"
	tcPubKeyAlias2    = "test_alias2"
	tcServiceName1    = "test_service1"
	tcServiceName2    = "test_service2"
	tcMessageSize     = (1 << 20)
	tcWorkSize        = 22
	tcFetchTimeout    = 5000
	tcQueuePeriod     = 1000
	tcQBPConsumers    = 5
	tcPowParallel     = 8
)

var (
	tgAdapters = []string{
		"test_connect1",
		"test_connect2",
	}
	tgPubKeys = map[string]string{
		tcPubKeyAlias1: tgPubKey1.ToString(),
		tcPubKeyAlias2: tgPubKey2.ToString(),
	}
	tgServices = map[string]string{
		tcServiceName1: "test_address1",
		tcServiceName2: "test_address2",
	}
)

const (
	tcConfigTemplate = `settings:
  message_size_bytes: %d
  work_size_bits: %d
  fetch_timeout_ms: %d
  queue_period_ms: %d
  network_key: %s
  qbp_consumers: %d
  pow_parallel: %d
logging:
  - info
  - erro
address:
  external: %s
  internal: %s
endpoints:
  - %s
  - %s
friends:
  %s: %s
  %s: %s
services:
  %s: %s
  %s: %s
`
)

func testNewConfigString() string {
	return fmt.Sprintf(
		tcConfigTemplate,
		tcMessageSize,
		tcWorkSize,
		tcFetchTimeout,
		tcQueuePeriod,
		tcNetwork,
		tcQBPConsumers,
		tcPowParallel,
		tcAddressExternal,
		tcAddressInternal,
		tgAdapters[0],
		tgAdapters[1],
		tcPubKeyAlias1,
		tgPubKeys[tcPubKeyAlias1],
		tcPubKeyAlias2,
		tgPubKeys[tcPubKeyAlias2],
		tcServiceName1,
		tgServices[tcServiceName1],
		tcServiceName2,
		tgServices[tcServiceName2],
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
		return errors.New("success load config on non exist file")
	}

	if err := os.WriteFile(configFile, []byte("abc"), 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid structure")
	}

	cfg1Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "settings", "settings_v2"))
	if err := os.WriteFile(configFile, cfg1Bytes, 0o600); err != nil {
		return err
	}

	cfg2Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "PubKey", "PubKey_v2"))
	if err := os.WriteFile(configFile, cfg2Bytes, 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (friends)")
	}

	cfg3Bytes := []byte(strings.ReplaceAll(testNewConfigString(), "erro", "erro_v2"))
	if err := os.WriteFile(configFile, cfg3Bytes, 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (logging)")
	}

	pubKey1 := tgPubKeys[tcPubKeyAlias1]
	pubKey2 := tgPubKeys[tcPubKeyAlias2]

	cfg4Bytes := []byte(strings.ReplaceAll(testNewConfigString(), pubKey1, pubKey2))
	if err := os.WriteFile(configFile, cfg4Bytes, 0o600); err != nil {
		return err
	}

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with invalid fields (duplicate publc keys)")
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

	if cfg.GetSettings().GetFetchTimeout() != time.Duration(tcFetchTimeout)*time.Millisecond {
		t.Fatal("settings fetch timeout is invalid")
	}

	if cfg.GetSettings().GetQueuePeriod() != time.Duration(tcQueuePeriod)*time.Millisecond {
		t.Fatal("settings queue period is invalid")
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

	if cfg.GetSettings().GetQBPConsumers() != tcQBPConsumers {
		t.Fatal("qbp_consumers is invalid")
	}

	if cfg.GetSettings().GetPowParallel() != tcPowParallel {
		t.Fatal("pow_parallel is invalid")
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
		if v != tgAdapters[i] {
			t.Fatalf("connection '%d' is invalid", i)
		}
	}

	for k, v := range tgServices {
		v1, ok := cfg.GetService(k)
		if !ok {
			t.Fatalf("service undefined '%s'", k)
		}
		if v != v1 {
			t.Fatalf("service host is invalid '%s'", v1)
		}
	}

	for name, pubStr := range tgPubKeys {
		v1 := cfg.GetFriends()[name]
		pubKey := asymmetric.LoadPubKey(pubStr)
		if pubKey.ToString() != v1.ToString() {
			t.Fatalf("public key is invalid '%s'", v1)
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

	if len(cfg.GetFriends()) == 0 {
		t.Fatal("list of friends should be is not nil for tests")
	}

	wrapper := NewWrapper(cfg)
	_ = wrapper.GetEditor().UpdateFriends(nil)

	if len(cfg.GetFriends()) != 0 {
		t.Fatal("friends is not nil for current config")
	}

	cfg, err = LoadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.GetFriends()) != 0 {
		t.Fatal("friends is not nil for loaded config")
	}
}
