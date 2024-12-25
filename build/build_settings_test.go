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

	settings.FQueueProblem.FMainPoolCap = 64
	settings.FQueueProblem.FRandPoolCap = 64
	settings.FQueueProblem.FConsumersCap = 1
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
	settings.FNetworkConnection.FReadTimeoutMS = 5_000
	settings.FNetworkConnection.FWriteTimeoutMS = 5_000
	settings.FNetworkConnection.FWaitTimeoutMS = 5_000_000
	if err := settings.validate(); err != nil {
		t.Error(err)
		return
	}

	if GSettings.FProtoMask.FNetwork != 0x5f67705f {
		t.Error(`GSettings.ProtoMask.Network != 0x5f67705f`)
		return
	}
	if GSettings.FProtoMask.FService != 0x5f686c5f {
		t.Error(`GGSettings.ProtoMask.Service != 0x5f686c5f`)
		return
	}
	if GSettings.FQueueProblem.FMainPoolCap != 256 {
		t.Error(`GSettings.QueueCapacity.FMainPoolCap != 256`)
		return
	}
	if GSettings.FQueueProblem.FRandPoolCap != 32 {
		t.Error(`GSettings.QueueCapacity.FRandPoolCap != 32`)
		return
	}
	if GSettings.FQueueProblem.FConsumersCap != 1 {
		t.Error(`GSettings.QueueCapacity.FConsumersCap != 1`)
		return
	}
	if GSettings.FNetworkManager.FCacheHashesCap != 2048 {
		t.Error(`GSettings.NetworkManager.CacheHashesCap != 2048`)
		return
	}
	if GSettings.FNetworkManager.FConnectsLimiter != 256 {
		t.Error(`GSettings.NetworkManager.ConnectsLimiter != 256`)
		return
	}
	if GSettings.FNetworkManager.FKeeperPeriodMS != 10_000 {
		t.Error(`GSettings.NetworkManager.KeeperPeriodMS != 10_000`)
		return
	}
	if GSettings.FNetworkConnection.FDialTimeoutMS != 5_000 {
		t.Error(`GSettings.NetworkConnection.DialTimeoutMS != 5_000`)
		return
	}
	if GSettings.FNetworkConnection.FReadTimeoutMS != 5_000 {
		t.Error(`GSettings.NetworkConnection.ReadTimeoutMS != 5_000`)
		return
	}
	if GSettings.FNetworkConnection.FWriteTimeoutMS != 5_000 {
		t.Error(`GSettings.NetworkConnection.WriteTimeoutMS != 5_000`)
		return
	}
	if GSettings.FNetworkConnection.FWaitTimeoutMS != 3_600_000 {
		t.Error(`GSettings.NetworkConnection.WaitTimeoutMS != 3_600_000`)
		return
	}
	switch {
	case GSettings.GetWaitTimeout() != time.Duration(GSettings.FNetworkConnection.FWaitTimeoutMS)*time.Millisecond:
		fallthrough
	case GSettings.GetDialTimeout() != time.Duration(GSettings.FNetworkConnection.FDialTimeoutMS)*time.Millisecond:
		fallthrough
	case GSettings.GetReadTimeout() != time.Duration(GSettings.FNetworkConnection.FReadTimeoutMS)*time.Millisecond:
		fallthrough
	case GSettings.GetWriteTimeout() != time.Duration(GSettings.FNetworkConnection.FWriteTimeoutMS)*time.Millisecond:
		fallthrough
	case GSettings.GetKeeperPeriod() != time.Duration(GSettings.FNetworkManager.FKeeperPeriodMS)*time.Millisecond:
		t.Error("Get methods (settings) is not valid")
	}
}
