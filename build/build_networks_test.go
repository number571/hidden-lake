package build

import (
	_ "embed"
	"testing"
)

func TestHiddenLakeNetworks(t *testing.T) {
	t.Parallel()

	network := SNetwork{}
	if err := network.validate(); err == nil {
		t.Error("success validate with invalid message_size_bytes")
		return
	}

	network.FMessageSizeBytes = 8192
	if err := network.validate(); err != nil {
		t.Error("error validate with exist message_size_bytes")
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

	if network.FWorkSizeBits != 0 {
		t.Error("network.FWorkSizeBits != 0")
		return
	}

	networks := GetNetworks()
	newNetwork := networks[CDefaultNetwork]

	newNetworkKey := "new_network"
	neMessageSize := uint64(9_123)

	if _, ok := networks[newNetworkKey]; ok {
		t.Error("new network key already exist")
		return
	}
	if newNetwork.FMessageSizeBytes == neMessageSize {
		t.Error("new set value already equal")
		return
	}

	newNetwork.FMessageSizeBytes = neMessageSize
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
	if gotNetwork.FMessageSizeBytes != neMessageSize {
		t.Error("new set value not saved")
		return
	}

	if err := SetNetworks(map[string]SNetwork{}); err == nil {
		t.Error("success set networks without default")
		return
	}
	if err := SetNetworks(map[string]SNetwork{CDefaultNetwork: SNetwork{}}); err == nil {
		t.Error("success set incorrect network")
		return
	}
}
