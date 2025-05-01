package build

import (
	_ "embed"
	"testing"
	"time"
)

func TestHiddenLakeNetworks(t *testing.T) {
	t.Parallel()

	network := SNetwork{}
	if err := network.validate(); err == nil {
		t.Error("success validate with invalid message_size_bytes")
		return
	}

	network.FMessageSizeBytes = 8192
	if err := network.validate(); err == nil {
		t.Error("success validate with invalid fetch_timeout_ms")
		return
	}

	network.FFetchTimeoutMS = 60_000
	if err := network.validate(); err == nil {
		t.Error("success validate with invalid queue_period_ms")
		return
	}

	network.FQueuePeriodMS = 5_000
	if err := network.validate(); err != nil {
		t.Error(err)
		return
	}

	network.FConnections = []string{"127.0.0.1:8080"}
	if err := network.validate(); err == nil {
		t.Error("success validate with invalid connections (1)")
		return
	}

	network.FConnections = []string{"127.0.0.1"}
	if err := network.validate(); err == nil {
		t.Error("success validate with invalid connections (2)")
		return
	}

	network, ok := gNetworks[CDefaultNetwork]
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

	networks := GetNetworks()
	newNetwork := networks[CDefaultNetwork]

	newNetworkKey := "new_network"
	newFetchTimeout := uint64(2_123)

	if _, ok := networks[newNetworkKey]; ok {
		t.Error("new network key already exist")
		return
	}
	if newNetwork.FFetchTimeoutMS == newFetchTimeout {
		t.Error("new set value already equal")
		return
	}

	newNetwork.FFetchTimeoutMS = newFetchTimeout
	networks[newNetworkKey] = newNetwork
	if err := SetNetworks(networks); err != nil {
		t.Error(err)
		return
	}

	gotNetwork, ok := GetNetwork(newNetworkKey)
	if !ok {
		t.Error("new set network key not saved")
		return
	}
	if gotNetwork.FFetchTimeoutMS != newFetchTimeout {
		t.Error("new set value not saved")
		return
	}
}
