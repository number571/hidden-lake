package name

import "testing"

func TestServiceName(t *testing.T) {
	serviceName := LoadAppName("hidden-lake-kernel")
	if serviceName.Format() != "Hidden Lake Kernel" {
		t.Fatal("invalid format name")
	}
	if serviceName.Short() != "HLK" {
		t.Fatal("invalid short name")
	}

	serviceWithSubName := LoadAppName("hidden-lake-adapters=common")
	if serviceWithSubName.Format() != "Hidden Lake Adapters = Common" {
		t.Fatal("invalid format name (with subname)")
	}
	if serviceWithSubName.Short() != "HLA=common" {
		t.Fatal("invalid short name (with subname)")
	}
}
