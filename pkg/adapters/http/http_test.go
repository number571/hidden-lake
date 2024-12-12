package http

import (
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Parallel()

	path := "/path"
	handler := NewHandler(path, func(_ http.ResponseWriter, _ *http.Request) {})
	if handler.GetPath() != path {
		t.Error("path is invalid")
		return
	}
	_ = handler.GetFunc()
}

func TestSettings(t *testing.T) {
	t.Parallel()

	sett := NewSettings(nil)
	if sett.GetAdapterSettings() == nil {
		t.Error("invalid adapter settings")
		return
	}
}
