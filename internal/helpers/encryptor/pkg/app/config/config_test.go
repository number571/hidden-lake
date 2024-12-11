package config

import (
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

var (
	tgPubKeys = map[string]string{
		tcPubKeyAlias1: tgPubKey1.ToString(),
		tcPubKeyAlias2: tgPubKey2.ToString(),
	}
)

const (
	tcConfigFile   = "config_test.txt"
	tcLogging      = true
	tcNetwork      = "test_network"
	tcAddress1     = "test_address1"
	tcAddress2     = "test_address2"
	tcMessageSize  = (1 << 20)
	tcWorkSize     = 22
	tcPubKeyAlias1 = "test_alias1"
	tcPubKeyAlias2 = "test_alias2"
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
			FMessageSizeBytes: tcMessageSize,
			FWorkSizeBits:     tcWorkSize,
			FNetworkKey:       tcNetwork,
		},
		FLogging: []string{"info", "erro"},
		FAddress: &SAddress{
			FInternal: tcAddress1,
			FPPROF:    tcAddress2,
		},
		FFriends: map[string]string{
			tcPubKeyAlias1: tgPubKey1.ToString(),
			tcPubKeyAlias2: tgPubKey2.ToString(),
		},
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

	if cfg.GetSettings().GetWorkSizeBits() != tcWorkSize {
		t.Error("settings work size is invalid")
		return
	}

	if cfg.GetSettings().GetMessageSizeBytes() != tcMessageSize {
		t.Error("settings message size is invalid")
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

	if cfg.GetAddress().GetInternal() != tcAddress1 {
		t.Error("address http is invalid")
		return
	}

	if cfg.GetAddress().GetPPROF() != tcAddress2 {
		t.Error("address pprof is invalid")
		return
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
