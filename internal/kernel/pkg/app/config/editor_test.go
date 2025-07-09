package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

var (
	tgPubKey1 = asymmetric.NewPrivKey().GetPubKey()
	tgPubKey2 = asymmetric.NewPrivKey().GetPubKey()
)

var (
	// tgNewConnections = []string{"a", "b", "c", "b"}
	tgNewFriends = map[string]asymmetric.IPubKey{
		"a": tgPubKey1,
		"b": tgPubKey2,
	}

	// duplicated public keys
	tgNewIncorrect2Friends = map[string]asymmetric.IPubKey{
		"a": tgPubKey1,
		"b": tgPubKey1,
	}
)

type tsConfig struct{}

var (
	_ IConfig = &tsConfig{}
)

func (p *tsConfig) GetSettings() IConfigSettings              { return nil }
func (p *tsConfig) GetLogging() logger.ILogging               { return nil }
func (p *tsConfig) GetShare() bool                            { return false }
func (p *tsConfig) GetAddress() IAddress                      { return nil }
func (p *tsConfig) GetNetworkKey() string                     { return "" }
func (p *tsConfig) GetEndpoints() []string                    { return nil }
func (p *tsConfig) GetFriends() map[string]asymmetric.IPubKey { return nil }
func (p *tsConfig) GetService(_ string) (string, bool)        { return "", false }

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

	config := wrapper.GetConfig()
	editor := wrapper.GetEditor()

	beforeFriends := config.GetFriends()

	if err := editor.UpdateFriends(tgNewFriends); err != nil {
		t.Fatal(err)
	}
	afterFriends := config.GetFriends()
	if len(afterFriends) != 2 {
		t.Fatal("failed deduplicate public keys (friends)")
	}
	for af := range afterFriends {
		if _, ok := beforeFriends[af]; ok {
			t.Fatal("beforeFriends == afterFriends")
		}
	}
	for nf := range tgNewFriends {
		if _, ok := afterFriends[nf]; !ok {
			t.Fatal("afterFriends != tgNewFriends")
		}
	}

	if err := editor.UpdateFriends(tgNewIncorrect2Friends); err == nil {
		t.Fatal("success update friends with duplicates")
	}
}

func TestIncorrectFilepathEditor(t *testing.T) {
	t.Parallel()

	configFile := fmt.Sprintf(tcConfigFileTemplate, 5)
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

	if err := editor.UpdateFriends(tgNewFriends); err == nil {
		t.Fatal("success update friends with incorrect filepath")
	}
}
