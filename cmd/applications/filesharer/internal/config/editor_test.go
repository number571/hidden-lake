package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/hidden-lake/internal/language"
	logger "github.com/number571/hidden-lake/internal/logger/std"
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
func (p *tsConfig) GetShare() bool                   { return false }
func (p *tsConfig) GetAddress() IAddress             { return nil }
func (p *tsConfig) GetNetworkKey() string            { return "" }
func (p *tsConfig) GetConnection() string            { return "" }
func (p *tsConfig) GetSecretKeys() map[string]string { return nil }
func (p *tsConfig) GetStoragePath() string           { return "" }

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

	config.fFilepath = random.NewCSPRNG().GetString(32)

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
