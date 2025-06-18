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
			t.Fatal("nothing panics")
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
		t.Fatal("success upload invalid file")
	}
}

func TestGetMessage(t *testing.T) {
	t.Parallel()

	if _, err := GetMessage(wrapText(""), ""); err == nil {
		t.Fatal("success get void message")
	}
	if _, err := GetMessage(wrapFile("", []byte{}), ""); err == nil {
		t.Fatal("success get void file")
	}
}
