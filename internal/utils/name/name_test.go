package name

import "testing"

func TestServiceName(t *testing.T) {
	serviceName := LoadServiceName("hidden-lake-service")
	if serviceName.Format() != "Hidden Lake Service" {
		t.Fatal("invalid format name")
	}
	if serviceName.Short() != "HLS" {
		t.Fatal("invalid short name")
	}

	serviceWithSubName := LoadServiceName("hidden-lake-adapters=common")
	if serviceWithSubName.Format() != "Hidden Lake Adapters = Common" {
		t.Fatal("invalid format name (with subname)")
	}
	if serviceWithSubName.Short() != "HLA=common" {
		t.Fatal("invalid short name (with subname)")
	}
}
