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
func (p *tsConfig) GetShare() bool                   { return false }
func (p *tsConfig) GetAddress() string               { return "" }
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
			t.Fatal("nothing panics")
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
	defer func() { _ = os.Remove(configFile) }()

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}

	wrapper := NewWrapper(cfg)
	editor := wrapper.GetEditor()

	res, err := language.ToILanguage("RUS")
	if err != nil {
		t.Fatal(err)
	}
	if err := editor.UpdateLanguage(res); err != nil {
		t.Fatal(err)
	}
}

func TestIncorrectFilepathEditor(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 3)
	defer func() { _ = os.Remove(configFile) }()

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}

	wrapper := NewWrapper(cfg)

	config := wrapper.GetConfig().(*SConfig)
	editor := wrapper.GetEditor()

	config.fFilepath = random.NewRandom().GetString(32)

	res, err := language.ToILanguage("RUS")
	if err != nil {
		t.Fatal(err)
	}
	if err := editor.UpdateLanguage(res); err == nil {
		t.Fatal("success update network key with incorrect filepath")
	}
}
