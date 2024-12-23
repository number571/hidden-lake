package config

import "testing"

func TestConfig(t *testing.T) {
	v := uint64(10)
	c := &SConfigSettings{
		FPayloadSizeBytes: v,
	}
	if c.GetPayloadSizeBytes() != v {
		t.Error("payload size bytes != v")
		return
	}
}
