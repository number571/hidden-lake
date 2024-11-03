package handler

import (
	"encoding/base64"
	"testing"

	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
)

const (
	tcFile        = "file.txt"
	tcFileEscaped = "file.txt&lt;b&gt;some&lt;/b&gt;.txt"
	tcText        = "<a href='xxx'>hello</a>"
	tcTextEscaped = "&lt;a href=&#39;xxx&#39;&gt;hello&lt;/a&gt;"
)

func TestDataTypePanic(t *testing.T) {
	t.Parallel()

	testUnwrapTextPanic(t)
	testUnwrapFilePanic(t)
}

func testUnwrapTextPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = unwrapText([]byte{}, false)
}

func testUnwrapFilePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_, _ = unwrapFile([]byte{}, false)
}

func TestDataType(t *testing.T) {
	t.Parallel()

	if isText([]byte{}) {
		t.Error("isText([]byte{}) = ok")
		return
	}
	if isFile([]byte{}) {
		t.Error("isFile([]byte{}) = ok")
		return
	}

	wt := wrapText(tcText)
	if !isText(wt) {
		t.Error("wrapText: !isText(wt)")
		return
	}
	if isFile(wt) {
		t.Error("wrapText: isFile(wt)")
		return
	}
	if wt[0] != hlm_settings.CIsText {
		t.Error("wrapText:  wt[0] != hlm_settings.CIsText")
		return
	}
	if unwrapText(wt, false) != tcText {
		t.Error("wrapText: unwrapText(wt, false) != tcText")
		return
	}
	if unwrapText(wt, true) != tcTextEscaped {
		t.Error("wrapText: unwrapText(wt, true) != tcTextEscaped")
		return
	}

	wf := wrapFile(tcFile, []byte(tcText))
	if !isFile(wf) {
		t.Error("wrapFile: !isFile(wf)")
		return
	}
	if isText(wf) {
		t.Error("wrapFile: isText(wf)")
		return
	}
	if wf[0] != hlm_settings.CIsFile || wf[len(tcFile)+1] != hlm_settings.CIsFile {
		t.Error("wrapFile: wf[0] != hlm_settings.CIsFile || wf[len(tcFile)+1] != hlm_settings.CIsFile")
		return
	}
	if f, b := unwrapFile(wf, false); f != tcFile || b != base64.StdEncoding.EncodeToString([]byte(tcText)) {
		t.Error("wrapText: f, b := unwrapFile(wf, false); f != tcFile || b != base64.StdEncoding.EncodeToString([]byte(tcText))")
		return
	}

	wfx := wrapFile(tcFile+"<b>some</b>"+".txt", []byte(tcText))
	if f, b := unwrapFile(wfx, true); f != tcFileEscaped || b != base64.StdEncoding.EncodeToString([]byte(tcText)) {
		t.Error("wrapText: f, b := unwrapFile(wf, true); f != tcFileEscaped || b != base64.StdEncoding.EncodeToString([]byte(tcText))")
		return
	}

	if f, b := unwrapFile([]byte{1}, false); f != "" || b != "" {
		t.Error(`wrapFile: f, b := unwrapFile([]byte{1}, false); f != "" || b != ""`)
		return
	}
	wf2 := wrapFile(tcFile+"\x01"+".txt", []byte(tcText))
	if f, b := unwrapFile(wf2, false); f != "" || b != "" {
		t.Error(`wrapFile: f, b := unwrapFile(wf2, false); f != "" || b != ""`)
		return
	}
	wf3 := wrapFile(tcFile, []byte{})
	if f, b := unwrapFile(wf3, false); f != "" || b != "" {
		t.Error(`wrapFile: f, b := unwrapFile(wf3, false); f != "" || b != ""`)
		return
	}
}
