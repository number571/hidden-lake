package build

import (
	"testing"

	"github.com/number571/hidden-lake/build"
)

func TestFailedSetBuildByPath(t *testing.T) {
	t.Parallel()

	if _, err := SetBuildByPath("testdata/failed/1"); err == nil {
		t.Fatal("success set build with invalid settings")
	}
	if _, err := SetBuildByPath("testdata/failed/2"); err == nil {
		t.Fatal("success set build with invalid networks")
	}
	if _, err := SetBuildByPath("testdata/failed/3"); err == nil {
		t.Fatal("success set build with invalid settings (yaml)")
	}
	if _, err := SetBuildByPath("testdata/failed/4"); err == nil {
		t.Fatal("success set build with invalid networks (yaml)")
	}
}

func TestSuccessSetBuildByPath(t *testing.T) {
	t.Parallel()

	if _, err := SetBuildByPath("testdata/success"); err != nil {
		t.Fatal(err)
	}

	settings := build.GetSettings()
	if settings.FProtoMask.FNetwork != 0x01 {
		t.Fatal("settings are not saved")
	}

	networks := build.GetNetworks()
	testNetwork, ok := networks["__testdata__"]
	if !ok || testNetwork.FMessageSizeBytes != 4097 {
		t.Fatal("networks are not saved")
	}

	if oks, err := SetBuildByPath("__not_found_path"); err != nil || oks[0] || oks[1] {
		t.Fatal("success build not found path")
	}

	settings = build.GetSettings()
	if settings.FProtoMask.FNetwork != 0x01 {
		t.Fatal("settings are rewrites success with not found path")
	}

	networks = build.GetNetworks()
	testNetwork, ok = networks["__testdata__"]
	if !ok || testNetwork.FMessageSizeBytes != 4097 {
		t.Fatal("networks are rewrites success with not found path")
	}
}
