package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/hidden-lake/internal/utils/language"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

const (
	tcConfigFileTemplate = "config_test_%d.txt"
)

type tsConfig struct{}

var (
	_ IConfig = &tsConfig{}
)

func (p *tsConfig) GetSettings() IConfigSettings     { return nil }
func (p *tsConfig) GetLanguage() language.ILanguage  { return 0 }
func (p *tsConfig) GetLogging() logger.ILogging      { return nil }
func (p *tsConfig) GetShareEnabled() bool            { return false }
func (p *tsConfig) GetAddress() IAddress             { return nil }
func (p *tsConfig) GetNetworkKey() string            { return "" }
func (p *tsConfig) GetConnection() string            { return "" }
func (p *tsConfig) GetStorageKey() string            { return "" }
func (p *tsConfig) GetSecretKeys() map[string]string { return nil }
func (p *tsConfig) GetChannels() []string            { return nil }

func TestPanicEditor(t *testing.T) {
	t.Parallel()

	for i := 0; i < 2; i++ {
		testPanicEditor(t, i)
	}
}

func testPanicEditor(t *testing.T, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	switch n {
	case 0:
		_ = newEditor(nil)
	case 1:
		_ = newEditor(&tsConfig{})
	}
}

func TestEditor(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 4)
	defer os.Remove(configFile)

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	wrapper := NewWrapper(cfg)
	editor := wrapper.GetEditor()

	res, err := language.ToILanguage("RUS")
	if err != nil {
		t.Error(err)
		return
	}
	if err := editor.UpdateLanguage(res); err != nil {
		t.Error(err)
		return
	}

	channels := []string{"111", "222", "333", "222"}
	if err := editor.UpdateChannels(channels); err != nil {
		t.Error(err)
		return
	}
}

func TestIncorrectFilepathEditor(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 3)
	defer os.Remove(configFile)

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	wrapper := NewWrapper(cfg)

	config := wrapper.GetConfig().(*SConfig)
	editor := wrapper.GetEditor()

	config.fFilepath = random.NewRandom().GetString(32)

	res, err := language.ToILanguage("RUS")
	if err != nil {
		t.Error(err)
		return
	}
	if err := editor.UpdateLanguage(res); err == nil {
		t.Error("success update network key with incorrect filepath")
		return
	}
}
