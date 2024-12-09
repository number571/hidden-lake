// nolint: goerr113
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
	tcAddressPPROF    = "test_address_pprof"
	tcPubKeyAlias1    = "test_alias1"
	tcPubKeyAlias2    = "test_alias2"
	tcServiceName1    = "test_service1"
	tcServiceName2    = "test_service2"
	tcMessageSize     = (1 << 20)
	tcWorkSize        = 22
	tcFetchTimeout    = 5000
	tcQueuePeriod     = 1000
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
logging:
  - info
  - erro
address:
  external: %s
  internal: %s
  pprof: %s
adapters:
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
		tcAddressExternal,
		tcAddressInternal,
		tcAddressPPROF,
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
	defer os.Remove(config1File)

	cfg, err := LoadConfig(config1File)
	if err != nil {
		t.Error(err)
		return
	}

	if _, err := BuildConfig(config2File, &SConfig{}); err == nil {
		t.Error("success build config with void structure")
		return
	}

	if _, err := BuildConfig(config2File, cfg.(*SConfig)); err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(config2File)

	if _, err := BuildConfig(config2File, cfg.(*SConfig)); err == nil {
		t.Error("success build already exist config")
		return
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

	if _, err := LoadConfig(configFile); err == nil {
		return errors.New("success load config with required fields (settings)")
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
	defer os.Remove(configFile)

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

	if cfg.GetSettings().GetFetchTimeout() != time.Duration(tcFetchTimeout)*time.Millisecond {
		t.Error("settings fetch timeout is invalid")
		return
	}

	if cfg.GetSettings().GetQueuePeriod() != time.Duration(tcQueuePeriod)*time.Millisecond {
		t.Error("settings queue period is invalid")
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

	if cfg.GetAddress().GetPPROF() != tcAddressPPROF {
		t.Error("address_pprof is invalid")
		return
	}

	if len(cfg.GetAdapters()) != 2 {
		t.Error("len connections != 2")
		return
	}
	for i, v := range cfg.GetAdapters() {
		if v != tgAdapters[i] {
			t.Errorf("connection '%d' is invalid", i)
			return
		}
	}

	for k, v := range tgServices {
		v1, ok := cfg.GetService(k)
		if !ok {
			t.Errorf("service undefined '%s'", k)
			return
		}
		if v != v1 {
			t.Errorf("service host is invalid '%s'", v1)
			return
		}
	}

	for name, pubStr := range tgPubKeys {
		v1 := cfg.GetFriends()[name]
		pubKey := asymmetric.LoadPubKey(pubStr)
		if pubKey.ToString() != v1.ToString() {
			t.Errorf("public key is invalid '%s'", v1)
			return
		}
	}
}

func TestWrapper(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 1)

	testConfigDefaultInit(configFile)
	defer os.Remove(configFile)

	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	if len(cfg.GetFriends()) == 0 {
		t.Error("list of friends should be is not nil for tests")
		return
	}

	wrapper := NewWrapper(cfg)
	_ = wrapper.GetEditor().UpdateFriends(nil)

	if len(cfg.GetFriends()) != 0 {
		t.Error("friends is not nil for current config")
		return
	}

	cfg, err = LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	if len(cfg.GetFriends()) != 0 {
		t.Error("friends is not nil for loaded config")
		return
	}
}
