package appname

import "testing"

func TestServiceName(t *testing.T) {
	serviceName := LoadAppName("hidden-lake-kernel")
	if serviceName.Full() != "Hidden Lake Kernel" {
		t.Fatal("invalid full name")
	}
	if serviceName.Short() != "HLK" {
		t.Fatal("invalid short name")
	}

	serviceWithSubName := LoadAppName("hidden-lake-adapters=common")
	if serviceWithSubName.Full() != "Hidden Lake Adapters = Common" {
		t.Fatal("invalid full name (with subname)")
	}
	if serviceWithSubName.Short() != "HLA=common" {
		t.Fatal("invalid short name (with subname)")
	}
}
