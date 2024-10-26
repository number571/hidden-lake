package hiddenlake

import (
	_ "embed"
	"regexp"
	"strings"
	"testing"
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

	if strings.HasSuffix(CVersion, "~") {
		if match[0][1] != CVersion {
			t.Error("the versions do not match")
			return
		}
	} else {
		// current version is always previous version in the changelog
		if match[1][1] != CVersion {
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

	if len(GNetworks) == 0 {
		t.Error("len GNetworks = 0")
		return
	}

	network, ok := GNetworks["_test_network_"]
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

	if len(network.FConnections) == 0 {
		t.Error("len(network.FConnections) == 0")
		return
	}

	connection := network.FConnections[0]
	if connection.FHost != "xxx.xxx.xxx.xxx" {
		t.Error(`connection.FHost != "xxx.xxx.xxx.xxx"`)
		return
	}
	if len(connection.FPorts) != 2 {
		t.Error(`len(connection.FPorts) != 2`)
		return
	}
	if connection.FPorts[0] != 1 {
		t.Error(`connection.FPorts[0] != 1`)
		return
	}
	if connection.FPorts[1] != 2 {
		t.Error(`connection.FPorts[1] != 2`)
		return
	}
	if connection.FProvider != "_provider_" {
		t.Error(`connection.FProvider != "_provider_"`)
		return
	}
	if connection.FLocation != "_location_" {
		t.Error(`connection.FLocation != "_location_"`)
		return
	}
	if connection.FLogging != true {
		t.Error(`connection.FLogging != true`)
		return
	}
}
