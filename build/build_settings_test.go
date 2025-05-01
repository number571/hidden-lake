package build

import (
	"testing"
	"time"
)

func TestHiddenLakeSettings(t *testing.T) {
	t.Parallel()

	settings := SSettings{}
	if err := settings.validate(); err == nil {
		t.Error("success validate with invalid queue capacity")
		return
	}

	settings.FQueueBasedProblem.FMainPoolCap = 64
	settings.FQueueBasedProblem.FRandPoolCap = 64
	settings.FQueueBasedProblem.FQBPConsumers = 1
	settings.FQueueBasedProblem.FPowParallels = 1
	if err := settings.validate(); err == nil {
		t.Error("success validate with invalid network manager")
		return
	}

	settings.FNetworkManager.FCacheHashesCap = 2048
	settings.FNetworkManager.FConnectsLimiter = 128
	settings.FNetworkManager.FKeeperPeriodMS = 5_000
	if err := settings.validate(); err == nil {
		t.Error("success validate with invalid network connection")
		return
	}

	settings.FNetworkConnection.FDialTimeoutMS = 5_000
	settings.FNetworkConnection.FRecvTimeoutMS = 5_000
	settings.FNetworkConnection.FSendTimeoutMS = 5_000
	settings.FNetworkConnection.FWaitTimeoutMS = 5_000_000
	if err := settings.validate(); err != nil {
		t.Error(err)
		return
	}

	if gSettings.FProtoMask.FNetwork != 0x5f67705f {
		t.Error(`gSettings.ProtoMask.Network != 0x5f67705f`)
		return
	}
	if gSettings.FProtoMask.FService != 0x5f686c5f {
		t.Error(`GGSettings.ProtoMask.Service != 0x5f686c5f`)
		return
	}
	if gSettings.FQueueBasedProblem.FMainPoolCap != 256 {
		t.Error(`gSettings.QueueCapacity.FMainPoolCap != 256`)
		return
	}
	if gSettings.FQueueBasedProblem.FRandPoolCap != 32 {
		t.Error(`gSettings.QueueCapacity.FRandPoolCap != 32`)
		return
	}
	if gSettings.FNetworkManager.FCacheHashesCap != 2048 {
		t.Error(`gSettings.NetworkManager.CacheHashesCap != 2048`)
		return
	}
	if gSettings.FNetworkManager.FConnectsLimiter != 256 {
		t.Error(`gSettings.NetworkManager.ConnectsLimiter != 256`)
		return
	}
	if gSettings.FNetworkManager.FKeeperPeriodMS != 10_000 {
		t.Error(`gSettings.NetworkManager.KeeperPeriodMS != 10_000`)
		return
	}
	if gSettings.FNetworkConnection.FDialTimeoutMS != 5_000 {
		t.Error(`gSettings.NetworkConnection.DialTimeoutMS != 5_000`)
		return
	}
	if gSettings.FNetworkConnection.FRecvTimeoutMS != 5_000 {
		t.Error(`gSettings.NetworkConnection.FRecvTimeoutMS != 5_000`)
		return
	}
	if gSettings.FNetworkConnection.FSendTimeoutMS != 5_000 {
		t.Error(`gSettings.NetworkConnection.FSendTimeoutMS != 5_000`)
		return
	}
	if gSettings.FNetworkConnection.FWaitTimeoutMS != 3_600_000 {
		t.Error(`gSettings.NetworkConnection.WaitTimeoutMS != 3_600_000`)
		return
	}
	switch {
	case gSettings.GetWaitTimeout() != time.Duration(gSettings.FNetworkConnection.FWaitTimeoutMS)*time.Millisecond: //nolint:gosec
		fallthrough
	case gSettings.GetDialTimeout() != time.Duration(gSettings.FNetworkConnection.FDialTimeoutMS)*time.Millisecond: //nolint:gosec
		fallthrough
	case gSettings.GetRecvTimeout() != time.Duration(gSettings.FNetworkConnection.FRecvTimeoutMS)*time.Millisecond: //nolint:gosec
		fallthrough
	case gSettings.GetSendTimeout() != time.Duration(gSettings.FNetworkConnection.FSendTimeoutMS)*time.Millisecond: //nolint:gosec
		fallthrough
	case gSettings.GetKeeperPeriod() != time.Duration(gSettings.FNetworkManager.FKeeperPeriodMS)*time.Millisecond: //nolint:gosec
		t.Error("Get methods (settings) is not valid")
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
