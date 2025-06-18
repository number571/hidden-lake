package build

import (
	_ "embed"
	"testing"
)

func TestHiddenLakeNetworks(t *testing.T) {
	t.Parallel()

	network := SNetwork{}
	if err := network.validate(); err == nil {
		t.Fatal("success validate with invalid message_size_bytes")
	}

	network.FMessageSizeBytes = 8192
	if err := network.validate(); err != nil {
		t.Fatal("error validate with exist message_size_bytes")
	}

	network.FConnections = []string{"127.0.0.1:8080"}
	if err := network.validate(); err == nil {
		t.Fatal("success validate with invalid connections (1)")
	}

	network.FConnections = []string{"127.0.0.1"}
	if err := network.validate(); err == nil {
		t.Fatal("success validate with invalid connections (2)")
	}

	network, ok := gNetworks[CDefaultNetwork]
	if !ok {
		t.Fatal("not found network _test_network_")
	}

	if network.FMessageSizeBytes != 8192 {
		t.Fatal("network.FMessageSizeBytes != 8192")
	}

	if network.FWorkSizeBits != 0 {
		t.Fatal("network.FWorkSizeBits != 0")
	}

	networks := GetNetworks()
	newNetwork := networks[CDefaultNetwork]

	newNetworkKey := "new_network"
	neMessageSize := uint64(9_123)

	if _, ok := networks[newNetworkKey]; ok {
		t.Fatal("new network key already exist")
	}
	if newNetwork.FMessageSizeBytes == neMessageSize {
		t.Fatal("new set value already equal")
	}

	newNetwork.FMessageSizeBytes = neMessageSize
	networks[newNetworkKey] = newNetwork
	if err := SetNetworks(networks); err != nil {
		t.Fatal(err)
	}

	gotNetwork, ok := GetNetwork(newNetworkKey)
	if !ok {
		t.Fatal("new set network key not saved")
	}
	if gotNetwork.FMessageSizeBytes != neMessageSize {
		t.Fatal("new set value not saved")
	}

	if err := SetNetworks(map[string]SNetwork{}); err == nil {
		t.Fatal("success set networks without default")
	}
	if err := SetNetworks(map[string]SNetwork{CDefaultNetwork: SNetwork{}}); err == nil {
		t.Fatal("success set incorrect network")
	}
}
