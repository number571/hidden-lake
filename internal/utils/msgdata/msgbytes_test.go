package msgdata

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestPanicGetMessageBytes(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	formData := url.Values{
		"method": {"DELETE"},
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
	_, _ = GetMessageBytes(req)
}

func TestGetUploadFile(t *testing.T) {
	t.Parallel()

	formData := url.Values{"input_file": {""}}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(formData.Encode()))
	if _, _, err := getUploadFile(req); err == nil {
		t.Error("success upload invalid file")
		return
	}
}

func TestGetMessage(t *testing.T) {
	t.Parallel()

	if _, err := GetMessage(wrapText(""), ""); err == nil {
		t.Error("success get void message")
		return
	}
	if _, err := GetMessage(wrapFile("", []byte{}), ""); err == nil {
		t.Error("success get void file")
		return
	}
}
