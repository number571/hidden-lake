package build

import (
	"testing"

	"github.com/number571/hidden-lake/build"
)

func TestFailedSetBuildByPath(t *testing.T) {
	t.Parallel()

	if _, err := SetBuildByPath("testdata/failed/1"); err == nil {
		t.Error("success set build with invalid settings")
		return
	}
	if _, err := SetBuildByPath("testdata/failed/2"); err == nil {
		t.Error("success set build with invalid networks")
		return
	}
	if _, err := SetBuildByPath("testdata/failed/3"); err == nil {
		t.Error("success set build with invalid settings (yaml)")
		return
	}
	if _, err := SetBuildByPath("testdata/failed/4"); err == nil {
		t.Error("success set build with invalid networks (yaml)")
		return
	}
}

func TestSuccessSetBuildByPath(t *testing.T) {
	t.Parallel()

	if _, err := SetBuildByPath("testdata/success"); err != nil {
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

	if oks, err := SetBuildByPath("__not_found_path"); err != nil || oks[0] || oks[1] {
		t.Error("success build not found path")
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
