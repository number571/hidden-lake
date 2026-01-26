package adapters

import (
	"testing"

	"github.com/number571/hidden-lake/build"
)

func TestSettings(t *testing.T) {
	t.Parallel()

	sett := NewSettings(nil)
	defaultNetwork, _ := build.GetNetwork(build.CDefaultNetwork)

	if sett.GetMessageSizeBytes() != defaultNetwork.FMessageSizeBytes {
		t.Fatal("get invalid settings")
	}

	_ = NewSettingsByNetworkKey(build.CDefaultNetwork)
}

func TestPanicSettings(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	_ = NewSettingsByNetworkKey("__test_unknown__")
}
