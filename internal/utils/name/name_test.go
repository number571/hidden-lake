package name

import "testing"

func TestServiceName(t *testing.T) {
	serviceName := LoadServiceName("hidden-lake-service")
	if serviceName.Format() != "Hidden Lake Service" {
		t.Error("invalid format name")
		return
	}
	if serviceName.Short() != "HLS" {
		t.Error("invalid short name")
		return
	}

	serviceWithSubName := LoadServiceName("hidden-lake-adapters=common")
	if serviceWithSubName.Format() != "Hidden Lake Adapters = Common" {
		t.Error("invalid format name (with subname)")
		return
	}
	if serviceWithSubName.Short() != "HLA=common" {
		t.Error("invalid short name (with subname)")
		return
	}
}
