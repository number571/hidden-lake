package sealed

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/hidden-lake/build"
)

var (
	tgAPIToken  = strings.Repeat("a1", 32)
	tgPushToken = strings.Repeat("b2", 32)
	tgSeenKey   = strings.Repeat("c3", 32)
)

func testBootMap() map[string]any {
	return map[string]any{
		"priv_key":     tgPrivKey1.ToString(),
		"network_key":  tcNetworkKey,
		"friends":      map[string]string{"bob": tgPrivKey2.GetPubKey().ToString()},
		"endpoints":    []string{tcRelay1, tcRelay2},
		"service_addr": tcServiceAddr,
		"api_token":    tgAPIToken,
		"push_token":   tgPushToken,
		"seen_db_path": "seen.db",
		"seen_key":     tgSeenKey,
	}
}

func testBootLine(t *testing.T, pMap map[string]any) []byte {
	t.Helper()

	b, err := json.Marshal(pMap)
	if err != nil {
		t.Fatal(err)
	}
	return append(b, '\n')
}

func TestParseBoot(t *testing.T) {
	t.Parallel()

	boot, err := ParseBoot(testBootLine(t, testBootMap()))
	if err != nil {
		t.Fatal(err)
	}

	if boot.IsPrivKeyGenerated() {
		t.Fatal("priv key marked as generated")
	}
	if boot.GetPrivKey().ToString() != tgPrivKey1.ToString() {
		t.Fatal("invalid priv key")
	}
	if boot.GetAPIToken() != tgAPIToken || boot.GetPushToken() != tgPushToken {
		t.Fatal("invalid tokens")
	}
	if boot.GetSeenDBPath() != "seen.db" {
		t.Fatal("invalid seen db path")
	}
	seenKey, _ := hex.DecodeString(tgSeenKey)
	if !bytes.Equal(boot.GetSeenKey(), seenKey) {
		t.Fatal("invalid seen key")
	}

	cfg := boot.GetConfig()
	if cfg.GetSettings().GetNetworkKey() != tcNetworkKey {
		t.Fatal("invalid network key")
	}
	if addr, ok := cfg.GetService(tcMessenger); !ok || addr != tcServiceAddr {
		t.Fatal("invalid service addr")
	}
	if len(cfg.GetEndpoints()) != 2 {
		t.Fatal("invalid endpoints")
	}
	friend, ok := cfg.GetFriends()["bob"]
	if !ok || friend.ToString() != tgPrivKey2.GetPubKey().ToString() {
		t.Fatal("invalid friends")
	}
}

func TestParseBootGenerateKey(t *testing.T) {
	t.Parallel()

	m := testBootMap()
	m["priv_key"] = ""
	m["network_key"] = ""
	delete(m, "friends")
	delete(m, "endpoints")

	boot, err := ParseBoot(testBootLine(t, m))
	if err != nil {
		t.Fatal(err)
	}
	if !boot.IsPrivKeyGenerated() || boot.GetPrivKey() == nil {
		t.Fatal("priv key is not generated")
	}
	if boot.GetPrivKey().ToString() == tgPrivKey1.ToString() {
		t.Fatal("priv key is not new")
	}
	if len(boot.GetConfig().GetFriends()) != 0 || len(boot.GetConfig().GetEndpoints()) != 0 {
		t.Fatal("unexpected friends or endpoints")
	}
}

func TestParseBootInvalid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		field string
		value any
		err   error
	}{
		{"short api token", "api_token", strings.Repeat("a", 31), ErrInvalidAPIToken},
		{"empty api token", "api_token", "", ErrInvalidAPIToken},
		{"short push token", "push_token", strings.Repeat("b", 31), ErrInvalidPushToken},
		{"short seen key", "seen_key", strings.Repeat("c", 62), ErrInvalidSeenKey},
		{"long seen key", "seen_key", strings.Repeat("c", 66), ErrInvalidSeenKey},
		{"not hex seen key", "seen_key", strings.Repeat("zz", 32), ErrInvalidSeenKey},
		{"empty seen db path", "seen_db_path", "", ErrInvalidSeenDBPath},
		{"invalid service addr", "service_addr", "localhost", ErrInvalidServiceAddr},
		{"invalid priv key", "priv_key", "PrivKey{00}", ErrInvalidPrivKey},
		{"pub key as priv key", "priv_key", tgPrivKey1.GetPubKey().ToString(), ErrInvalidPrivKey},
		{"unknown network", "network_key", tcUnknownNetKey, ErrNetworkNotFound},
		{"invalid friend", "friends", map[string]string{"bob": "PubKey{00}"}, ErrInvalidParticipantKey},
		{"priv key as friend", "friends", map[string]string{"bob": tgPrivKey2.ToString()}, ErrInvalidParticipantKey},
		{"zeros priv key", "priv_key", "PrivKey{" + strings.Repeat("00", asymmetric.CPrivKeySize) + "}", ErrInvalidPrivKey},
		{"FF priv key", "priv_key", "PrivKey{" + strings.Repeat("FF", asymmetric.CPrivKeySize) + "}", ErrInvalidPrivKey},
		{"FF friend key", "friends", map[string]string{"bob": "PubKey{" + strings.Repeat("FF", asymmetric.CPubKeySize) + "}"}, ErrInvalidParticipantKey},
		{"63 hex api token", "api_token", strings.Repeat("a", 63), ErrInvalidAPIToken},
		{"65 hex api token", "api_token", strings.Repeat("a", 65), ErrInvalidAPIToken},
		{"66 hex api token", "api_token", strings.Repeat("a", 66), ErrInvalidAPIToken},
		{"not hex api token", "api_token", strings.Repeat("g", 64), ErrInvalidAPIToken},
		{"not hex push token", "push_token", strings.Repeat("x", 64), ErrInvalidPushToken},
		{"66 hex push token", "push_token", strings.Repeat("b", 66), ErrInvalidPushToken},
		{"hostname service addr", "service_addr", "localhost:9551", ErrInvalidServiceAddr},
		{"not loopback service addr", "service_addr", "192.168.1.10:9551", ErrInvalidServiceAddr},
		{"zero port service addr", "service_addr", "127.0.0.1:0", ErrInvalidServiceAddr},
		{"default network", "network_key", build.CDefaultNetwork, ErrUnsupportedNetwork},
		{"symmetric network", "network_key", "std-internal-meshtastic", ErrUnsupportedNetwork},
		{"empty alias", "friends", map[string]string{"  ": tgPrivKey2.GetPubKey().ToString()}, ErrInvalidAliasName},
		{"spaced alias", "friends", map[string]string{" bob": tgPrivKey2.GetPubKey().ToString()}, ErrInvalidAliasName},
		{"control alias", "friends", map[string]string{"bo\nb": tgPrivKey2.GetPubKey().ToString()}, ErrInvalidAliasName},
		{"del alias", "friends", map[string]string{"bo\x7fb": tgPrivKey2.GetPubKey().ToString()}, ErrInvalidAliasName},
		{"unicode control alias", "friends", map[string]string{"bo\u0085b": tgPrivKey2.GetPubKey().ToString()}, ErrInvalidAliasName},
		{"slash alias", "friends", map[string]string{"bo/b": tgPrivKey2.GetPubKey().ToString()}, ErrInvalidAliasName},
		{"no port endpoint", "endpoints", []string{"127.0.0.1"}, ErrInvalidEndpoint},
		{"empty host endpoint", "endpoints", []string{":9581"}, ErrInvalidEndpoint},
		{"empty port endpoint", "endpoints", []string{"127.0.0.1:"}, ErrInvalidEndpoint},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := testBootMap()
			m[tc.field] = tc.value
			if _, err := ParseBoot(testBootLine(t, m)); !errors.Is(err, tc.err) {
				t.Fatalf("expected %v, got %v", tc.err, err)
			}
		})
	}

	if _, err := ParseBoot([]byte("{not json\n")); !errors.Is(err, ErrDecodeBoot) {
		t.Fatal("invalid json accepted")
	}
	if _, err := ParseBoot([]byte(`{"api_token":1}`)); !errors.Is(err, ErrDecodeBoot) {
		t.Fatal("invalid json types accepted")
	}

	m := testBootMap()
	m["unknown_field"] = "x"
	if _, err := ParseBoot(testBootLine(t, m)); !errors.Is(err, ErrDecodeBoot) {
		t.Fatal("unknown field accepted")
	}
}

func TestParseBootBoundary(t *testing.T) {
	t.Parallel()

	m := testBootMap()
	m["api_token"] = strings.Repeat("aB", 32) // exactly 64 hex, mixed case
	m["push_token"] = strings.Repeat("09", 32)
	m["service_addr"] = "[::1]:9551"
	m["friends"] = map[string]string{"Боб Иванов": tgPrivKey2.GetPubKey().ToString()}
	m["endpoints"] = []string{"relay.example.com:9581", "[::1]:9582"}
	boot, err := ParseBoot(testBootLine(t, m))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := boot.GetConfig().GetFriends()["Боб Иванов"]; !ok {
		t.Fatal("unicode alias rejected")
	}

	// returned seen key is a copy
	boot.GetSeenKey()[0] ^= 0xFF
	seenKey, _ := hex.DecodeString(tgSeenKey)
	if !bytes.Equal(boot.GetSeenKey(), seenKey) {
		t.Fatal("seen key is not copied")
	}
}
