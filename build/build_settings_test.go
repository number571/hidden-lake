package build

import (
	"testing"
	"time"
)

func TestHiddenLakeSettings(t *testing.T) {
	t.Parallel()

	settings := SSettings{}
	if err := settings.validate(); err == nil {
		t.Fatal("success validate with invalid storage_manager")
	}

	settings.FStorageManager.FCacheHashesCap = 2048
	settings.FStorageManager.FQueuePoolCap = [2]uint64{256, 32}
	if err := settings.validate(); err == nil {
		t.Fatal("success validate with invalid network_manager")
	}

	settings.FNetworkManager.FConnNumLimit = 2048
	settings.FNetworkManager.FMessageBuffer = 256
	settings.FNetworkManager.FHttpReadTimeoutMS = 5000
	settings.FNetworkManager.FHttpHandleTimeoutMS = 5000
	settings.FNetworkManager.FHttpCallbackTimeoutMS = 3600000
	if err := settings.validate(); err != nil {
		t.Fatal(err)
	}

	if gSettings.FStorageManager.FCacheHashesCap != 2048 {
		t.Fatal(`gSettings.FStorageManager.CacheHashesCap != 2048`)
	}
	if gSettings.GetHttpReadTimeout() != time.Duration(5_000)*time.Millisecond {
		t.Fatal(`gSettings.GetHttpReadTimeout() != time.Duration(5_000)*time.Millisecond`)
	}
	if gSettings.GetHttpHandleTimeout() != time.Duration(5_000)*time.Millisecond {
		t.Fatal(`gSettings.GetHttpHandleTimeout() != time.Duration(5_000)*time.Millisecond`)
	}
	if gSettings.GetHttpCallbackTimeout() != time.Duration(3_600_000)*time.Millisecond {
		t.Fatal(`gSettings.GetHttpCallbackTimeout() != time.Duration(3_600_000)*time.Millisecond`)
	}
	if gSettings.FStorageManager.FQueuePoolCap[0] != 256 {
		t.Fatal(`gSettings.FStorageManager.FQueuePoolCap[0] != 256`)
	}
	if gSettings.FStorageManager.FQueuePoolCap[1] != 32 {
		t.Fatal(`gSettings.FStorageManager.FQueuePoolCap[1] != 32`)
	}

	if err := SetSettings(SSettings{}); err == nil {
		t.Fatal("success set incorrect settings")
	}
}
