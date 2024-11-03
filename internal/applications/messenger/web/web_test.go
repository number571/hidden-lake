package web

import (
	"testing"
)

func TestWeb(t *testing.T) {
	t.Parallel()

	if !cUsedEmbedFS {
		t.Error("cUsedEmbedFS should be = true")
		return
	}
}
