package sealed

import (
	"errors"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/api/kernel/client/scheme"
)

const (
	tcServiceAddr   = "127.0.0.1:9551"
	tcMessenger     = "hls-messenger"
	tcNetworkKey    = "std-external-tcp"
	tcUnknownNetKey = "unknown-network-key"
	tcRelay1        = "127.0.0.1:9581"
	tcRelay2        = "127.0.0.1:9582"
)

var (
	tgPrivKey1 = asymmetric.NewPrivKey()
	tgPrivKey2 = asymmetric.NewPrivKey()
)

func testNewConfig(t *testing.T, pNetworkKey string) IConfig {
	t.Helper()

	cfg, err := NewConfig(
		pNetworkKey,
		tcServiceAddr,
		[]string{tcRelay1, tcRelay1},
		map[string]string{"alice": tgPrivKey1.GetPubKey().ToString()},
	)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestConfigFriends(t *testing.T) {
	t.Parallel()

	cfg := testNewConfig(t, "")

	friends := cfg.GetFriends()
	if len(friends) != 1 || friends["alice"].ToString() != tgPrivKey1.GetPubKey().ToString() {
		t.Fatal("invalid initial friends")
	}

	// the returned map is a copy
	delete(friends, "alice")
	if len(cfg.GetFriends()) != 1 {
		t.Fatal("friends map is not copied")
	}

	newFriends := map[string]layer2.IParticipantKey{
		"alice": tgPrivKey1.GetPubKey(),
		"bob":   tgPrivKey2.GetPubKey(),
	}
	if err := cfg.GetEditor().UpdateFriends(newFriends); err != nil {
		t.Fatal(err)
	}
	if len(cfg.GetConfig().GetFriends()) != 2 {
		t.Fatal("friends not updated")
	}

	dupFriends := map[string]layer2.IParticipantKey{
		"alice": tgPrivKey1.GetPubKey(),
		"carol": tgPrivKey1.GetPubKey(),
	}
	if err := cfg.UpdateFriends(dupFriends); !errors.Is(err, ErrDuplicateParticipantKey) {
		t.Fatal("duplicate participant key accepted")
	}
	if len(cfg.GetFriends()) != 2 {
		t.Fatal("friends changed after failed update")
	}
}

func TestConfigInvalidFriends(t *testing.T) {
	t.Parallel()

	_, err := NewConfig("", tcServiceAddr, nil, map[string]string{"alice": "PubKey{00}"})
	if !errors.Is(err, ErrInvalidParticipantKey) {
		t.Fatal("invalid participant key accepted")
	}

	pubKey := tgPrivKey1.GetPubKey().ToString()
	_, err = NewConfig("", tcServiceAddr, nil, map[string]string{"alice": pubKey, "bob": pubKey})
	if !errors.Is(err, ErrDuplicateParticipantKey) {
		t.Fatal("duplicate participant key accepted")
	}
}

func TestConfigService(t *testing.T) {
	t.Parallel()

	cfg := testNewConfig(t, "")

	addr, ok := cfg.GetService(tcMessenger)
	if !ok || addr != tcServiceAddr {
		t.Fatal("invalid messenger service")
	}
	if _, ok := cfg.GetService("hls-filesharer"); ok {
		t.Fatal("unexpected service")
	}
}

func TestConfigNetworkSettings(t *testing.T) {
	t.Parallel()

	network, ok := build.GetNetwork(tcNetworkKey)
	if !ok {
		t.Fatal("test network not found")
	}

	sett := testNewConfig(t, tcNetworkKey).GetSettings()
	if sett.GetNetworkKey() != tcNetworkKey {
		t.Fatal("invalid network key")
	}
	if sett.GetMessageSizeBytes() != network.FMessageSizeBytes {
		t.Fatal("invalid message size")
	}
	if sett.GetWorkSizeBits() != network.FWorkSizeBits {
		t.Fatal("invalid work size bits")
	}
	if uint64(sett.GetQueuePeriod().Milliseconds()) != network.FQueuePeriodMS { //nolint:gosec
		t.Fatal("invalid queue period")
	}
	if uint64(sett.GetFetchTimeout().Milliseconds()) != network.FFetchTimeoutMS { //nolint:gosec
		t.Fatal("invalid fetch timeout")
	}
	if sett.GetCryptoSchemeType() != scheme.GetCryptoSchemeType(network.FCryptoSchemeType) {
		t.Fatal("invalid crypto scheme type")
	}

	_, err := NewConfig(tcUnknownNetKey, tcServiceAddr, nil, nil)
	if !errors.Is(err, ErrNetworkNotFound) {
		t.Fatal("unknown network key accepted")
	}
}

func TestConfigDefaultSettings(t *testing.T) {
	t.Parallel()

	// same as hlk.yml created by default (without settings)
	cfg := testNewConfig(t, "")
	sett := cfg.GetSettings()
	if sett.GetNetworkKey() != "" ||
		sett.GetMessageSizeBytes() != 0 ||
		sett.GetWorkSizeBits() != 0 ||
		sett.GetQueuePeriod() != 0 ||
		sett.GetFetchTimeout() != 0 ||
		sett.GetQBPConsumers() != 0 ||
		sett.GetPowParallel() != 0 {
		t.Fatal("invalid default settings")
	}
	if sett.GetCryptoSchemeType() != scheme.CHybridScheme {
		t.Fatal("invalid default crypto scheme")
	}

	logging := cfg.GetLogging()
	// info is off by default: it would leak message-routing metadata to stderr.
	if logging.HasInfo() || !logging.HasWarn() || !logging.HasErro() {
		t.Fatal("invalid default logging")
	}
	if cfg.GetAddress() == nil {
		t.Fatal("address is nil")
	}
}

func TestConfigEndpoints(t *testing.T) {
	t.Parallel()

	cfg := testNewConfig(t, "")

	endpoints := cfg.GetEndpoints()
	if len(endpoints) != 1 || endpoints[0] != tcRelay1 {
		t.Fatal("initial endpoints are not deduplicated")
	}

	endpoints[0] = "changed"
	if cfg.GetEndpoints()[0] != tcRelay1 {
		t.Fatal("endpoints slice is not copied")
	}

	cfg.AddEndpoint(tcRelay2)
	cfg.AddEndpoint(tcRelay2)
	if len(cfg.GetEndpoints()) != 2 {
		t.Fatal("add endpoint is not deduplicated")
	}

	cfg.DelEndpoint(tcRelay1)
	cfg.DelEndpoint(tcRelay1)
	endpoints = cfg.GetEndpoints()
	if len(endpoints) != 1 || endpoints[0] != tcRelay2 {
		t.Fatal("del endpoint failed")
	}
}

func TestConfigUpdateFriendsAlias(t *testing.T) {
	t.Parallel()

	cfg := testNewConfig(t, "")
	for _, alias := range []string{"", " ", " bob", "bo\nb", "bo/b", "bo\u200eb"} {
		err := cfg.UpdateFriends(map[string]layer2.IParticipantKey{alias: tgPrivKey2.GetPubKey()})
		if !errors.Is(err, ErrInvalidAliasName) {
			t.Fatalf("invalid alias %q accepted", alias)
		}
	}
	if len(cfg.GetFriends()) != 1 {
		t.Fatal("friends changed after failed update")
	}
}
