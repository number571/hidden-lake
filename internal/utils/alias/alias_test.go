package alias

import (
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func TestGetAliasesList(t *testing.T) {
	t.Parallel()

	mapKeys := map[string]asymmetric.IPubKey{
		"AAA": asymmetric.NewPrivKey().GetPubKey(),
		"BBB": asymmetric.NewPrivKey().GetPubKey(),
		"CCC": asymmetric.NewPrivKey().GetPubKey(),
	}
	aliases := GetAliasesList(mapKeys)
	for _, alias := range aliases {
		if _, ok := mapKeys[alias]; !ok {
			t.Fatal("alias not found in map")
		}
	}
}

func TestGetAliasByPubKey(t *testing.T) {
	t.Parallel()

	bbbKey := asymmetric.NewPrivKey().GetPubKey()
	mapKeys := map[string]asymmetric.IPubKey{
		"AAA": asymmetric.NewPrivKey().GetPubKey(),
		"BBB": bbbKey,
		"CCC": asymmetric.NewPrivKey().GetPubKey(),
	}
	alias := GetAliasByPubKey(mapKeys, bbbKey)
	if alias != "BBB" {
		t.Fatal("get invalid alias name")
	}
	if GetAliasByPubKey(mapKeys, asymmetric.NewPrivKey().GetPubKey()) != "" {
		t.Fatal("get alias by unknown pub key")
	}
}
