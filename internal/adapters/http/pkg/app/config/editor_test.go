package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/random"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

var (
	tgNewConnections = []string{"a", "b"}
)

type tsConfig struct{}

var (
	_ IConfig = &tsConfig{}
)

func (p *tsConfig) GetSettings() IConfigSettings { return nil }
func (p *tsConfig) GetLogging() logger.ILogging  { return nil }
func (p *tsConfig) GetShare() bool               { return false }
func (p *tsConfig) GetAddress() IAddress         { return nil }
func (p *tsConfig) GetEndpoints() []string       { return nil }
func (p *tsConfig) GetConnections() []string     { return nil }

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
	defer func() { _ = os.Remove(configFile) }()

	testConfigDefaultInit(configFile)
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Error(err)
		return
	}

	wrapper := NewWrapper(cfg)

	config := wrapper.GetConfig()
	editor := wrapper.GetEditor()

	if err := editor.UpdateConnections(tgNewConnections); err != nil {
		t.Error(err)
		return
	}
	afterConnections := config.GetConnections()
	if len(afterConnections) != 2 {
		t.Error("failed deduplicate public keys (friends)")
		return
	}
	for i := range afterConnections {
		if tgNewConnections[i] != afterConnections[i] {
			t.Error("invalid new connections")
			return
		}
	}
}

func TestIncorrectFilepathEditor(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 5)
	defer func() { _ = os.Remove(configFile) }()

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

	if err := editor.UpdateConnections(tgNewConnections); err == nil {
		t.Error("success update friends with incorrect filepath")
		return
	}
}
