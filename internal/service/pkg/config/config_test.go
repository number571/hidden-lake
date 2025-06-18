package config

import "testing"

func TestConfig(t *testing.T) {
	v := uint64(10)
	c := &SConfigSettings{
		FPayloadSizeBytes: v,
	}
	if c.GetPayloadSizeBytes() != v {
		t.Fatal("payload size bytes != v")
	}
}
