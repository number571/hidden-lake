package hiddenlake

import (
	_ "embed"
	"regexp"
	"strings"
	"testing"
	"time"
)

var (
	//go:embed CHANGELOG.md
	tgCHANGELOG string
)

func TestHiddenLakeVersion(t *testing.T) {
	t.Parallel()

	re := regexp.MustCompile(`##\s+(v\d+\.\d+\.\d+~?)\s+`)
	match := re.FindAllStringSubmatch(tgCHANGELOG, -1)
	if len(match) < 2 {
		t.Error("versions not found")
		return
	}

	if strings.HasSuffix(GVersion, "~") {
		if match[0][1] != GVersion {
			t.Error("the versions do not match")
			return
		}
	} else {
		// current version is always previous version in the changelog
		if match[1][1] != GVersion {
			t.Error("the versions do not match")
			return
		}
	}

	if match[0][1] == match[1][1] {
		t.Error("the same versions inline")
		return
	}

	for i := 0; i < len(match); i++ {
		for j := i + 1; j < len(match)-1; j++ {
			if match[i][1] == match[j][1] {
				t.Errorf("found the same versions (i=%d, j=%d)", i, j)
				return
			}
		}
	}

	if strings.Count(tgCHANGELOG, "*??? ??, ????*") != 1 {
		t.Error("is there no new version or more than one new version?")
		return
	}
}

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

	if network.FWorkSizeBits != 22 {
		t.Error("network.FWorkSizeBits != 22")
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
}

func TestHiddenLakeSettings(t *testing.T) {
	t.Parallel()

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
		t.Error("Get methods is not valid")
	}
}
