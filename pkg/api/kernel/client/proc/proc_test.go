package proc

import (
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

func TestXxx(t *testing.T) {
	t.Parallel()

	sortList := []string{"Alice", "Bob", "Carol"}
	fMap := make(map[string]asymmetric.IPubKey, len(sortList))
	for _, v := range sortList {
		fMap[v] = asymmetric.NewPrivKey().GetPubKey()
	}

	list := FriendsMapToList(fMap)
	if len(list) != 3 {
		t.Fatal("len(list) != 3")
	}

	for i := range list {
		if list[i].FAliasName != sortList[i] {
			t.Fatal("sort failed")
		}
	}
}
