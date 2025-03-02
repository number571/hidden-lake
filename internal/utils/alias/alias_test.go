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
			t.Error("alias not found in map")
			return
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
		t.Error("get invalid alias name")
		return
	}
	if GetAliasByPubKey(mapKeys, asymmetric.NewPrivKey().GetPubKey()) != "" {
		t.Error("get alias by unknown pub key")
		return
	}
}
