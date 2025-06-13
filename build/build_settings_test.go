package build

import (
	"testing"
	"time"
)

func TestHiddenLakeSettings(t *testing.T) {
	t.Parallel()

	settings := SSettings{}
	if err := settings.validate(); err == nil {
		t.Error("success validate with invalid settings")
		return
	}

	settings.FStorageManager.FCacheHashesCap = 2048
	if err := settings.validate(); err == nil {
		t.Error("success validate with invalid network_manager")
		return
	}

	settings.FNetworkManager.FHttpReadTimeoutMS = 5000
	settings.FNetworkManager.FHttpHandleTimeoutMS = 5000
	settings.FNetworkManager.FHttpCallbackTimeoutMS = 3600000
	if err := settings.validate(); err == nil {
		t.Error("success validate with invalid queue_based_problem")
		return
	}

	settings.FQueueBasedProblem.FPoolCap = [2]uint64{256, 32}
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
	if gSettings.FStorageManager.FCacheHashesCap != 2048 {
		t.Error(`gSettings.FStorageManager.CacheHashesCap != 2048`)
		return
	}
	if gSettings.GetHttpReadTimeout() != time.Duration(5_000)*time.Millisecond {
		t.Error(`gSettings.GetHttpReadTimeout() != time.Duration(5_000)*time.Millisecond`)
		return
	}
	if gSettings.GetHttpHandleTimeout() != time.Duration(5_000)*time.Millisecond {
		t.Error(`gSettings.GetHttpHandleTimeout() != time.Duration(5_000)*time.Millisecond`)
		return
	}
	if gSettings.GetHttpCallbackTimeout() != time.Duration(3_600_000)*time.Millisecond {
		t.Error(`gSettings.GetHttpCallbackTimeout() != time.Duration(3_600_000)*time.Millisecond`)
		return
	}
	if gSettings.FQueueBasedProblem.FPoolCap[0] != 256 {
		t.Error(`gSettings.FQueueBasedProblem.FPoolCap[0] != 256`)
		return
	}
	if gSettings.FQueueBasedProblem.FPoolCap[1] != 32 {
		t.Error(`gSettings.FQueueBasedProblem.FPoolCap[1] != 32`)
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
