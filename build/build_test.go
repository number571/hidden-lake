package build

import (
	_ "embed"
	"testing"
	"time"
)

func TestHiddenLakeNetworks(t *testing.T) {
	t.Parallel()

	network, ok := GNetworks[CDefaultNetwork]
	if !ok {
		t.Error("not found network _test_network_")
		return
	}

	if network.FMessageSizeBytes != 8192 {
		t.Error("network.FMessageSizeBytes != 8192")
		return
	}

	if network.FFetchTimeoutMS != 60_000 {
		t.Error("network.FFetchTimeoutMS != 60_000")
		return
	}

	if network.FQueuePeriodMS != 5_000 {
		t.Error("network.FQueuePeriodMS != 5_000")
		return
	}

	if network.FWorkSizeBits != 0 {
		t.Error("network.FWorkSizeBits != 0")
		return
	}

	switch {
	case network.GetFetchTimeout() != time.Duration(60_000)*time.Millisecond:
		fallthrough
	case network.GetQueuePeriod() != time.Duration(5_000)*time.Millisecond:
		t.Error("Get methods (networks) is not valid")
	}
}

func TestHiddenLakeSettings(t *testing.T) {
	t.Parallel()

	settings := SSettings{}
	if err := settings.validate(); err == nil {
		t.Error("success validate with invalid queue capacity")
		return
	}

	settings.FQueueCapacity.FMain = 64
	settings.FQueueCapacity.FRand = 64
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
	if GSettings.FQueueCapacity.FMain != 256 {
		t.Error(`GSettings.QueueCapacity.Main != 256`)
		return
	}
	if GSettings.FQueueCapacity.FRand != 32 {
		t.Error(`GSettings.QueueCapacity.Rand != 32`)
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
