package build

import (
	"testing"

	"github.com/number571/hidden-lake/build"
)

func TestSetBuildByPath(t *testing.T) {
	t.Parallel()

	if err := SetBuildByPath("testdata"); err != nil {
		t.Error(err)
		return
	}

	settings := build.GetSettings()
	if settings.FProtoMask.FNetwork != 0x01 {
		t.Error("settings are not saved")
		return
	}

	networks := build.GetNetworks()
	testNetwork, ok := networks["__testdata__"]
	if !ok || testNetwork.FQueuePeriodMS != 1234 {
		t.Error("networks are not saved")
		return
	}

	if err := SetBuildByPath("__not_found_path"); err != nil {
		t.Error(err)
		return
	}

	settings = build.GetSettings()
	if settings.FProtoMask.FNetwork != 0x01 {
		t.Error("settings are rewrites success with not found path")
		return
	}

	networks = build.GetNetworks()
	testNetwork, ok = networks["__testdata__"]
	if !ok || testNetwork.FQueuePeriodMS != 1234 {
		t.Error("networks are rewrites success with not found path")
		return
	}
}
