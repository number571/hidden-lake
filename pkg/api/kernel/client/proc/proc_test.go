package proc

import (
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
)

func TestFriendsMapToList(t *testing.T) {
	t.Parallel()

	sortList := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	fMap := make(map[string]layer2.IParticipantKey, len(sortList))
	for _, v := range sortList {
		fMap[v] = asymmetric.NewPrivKey().GetPubKey()
	}

	list := FriendsMapToList(fMap)
	if len(list) != 10 {
		t.Fatal("len(list) != 10")
	}

	for i := range list {
		if list[i][0] != sortList[i] {
			t.Fatal("sort failed")
		}
	}
}
