package build

import (
	"testing"
)

func TestHiddenLakeSettings(t *testing.T) {
	t.Parallel()

	settings := SSettings{}
	if err := settings.validate(); err == nil {
		t.Error("success validate with invalid queue capacity")
		return
	}

	settings.FNetworkManager.FCacheHashesCap = 2048
	if err := settings.validate(); err != nil {
		t.Error(err)
		return
	}

	if gSettings.FProtoMask.FNetwork != 0x5f67705f {
		t.Error(`gSettings.ProtoMask.Network != 0x5f67705f`)
		return
	}
	if gSettings.FProtoMask.FService != 0x5f686c5f {
		t.Error(`gSettings.ProtoMask.Service != 0x5f686c5f`)
		return
	}
	if gSettings.FNetworkManager.FCacheHashesCap != 2048 {
		t.Error(`gSettings.NetworkManager.CacheHashesCap != 2048`)
		return
	}

	newSettings := GetSettings()

	newProtoMaskNetwork := uint32(0x1)
	if newSettings.FProtoMask.FNetwork == newProtoMaskNetwork {
		t.Error("new set value already equal")
		return
	}

	newSettings.FProtoMask.FNetwork = newProtoMaskNetwork
	if err := SetSettings(newSettings); err != nil {
		t.Error(err)
		return
	}

	if newSettings.FProtoMask.FNetwork != newProtoMaskNetwork {
		t.Error("new set value not saved")
		return
	}

	if err := SetSettings(SSettings{}); err == nil {
		t.Error("success set incorrect settings")
		return
	}
}
