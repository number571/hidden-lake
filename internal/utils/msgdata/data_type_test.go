package msgdata

import (
	"encoding/base64"
	"testing"
)

const (
	tcFile        = "file10.txt"
	tcFileEscaped = "file10.txt_b_some__b_.txt"
	tcText        = "<a href='xxx'>hello</a>"
	tcTextEscaped = "&lt;a href=&#39;xxx&#39;&gt;hello&lt;/a&gt;"
)

func TestDataType(t *testing.T) {
	t.Parallel()

	if text := unwrapText([]byte{}); text != "" {
		t.Fatal("unwrapText([]byte{}) = ok")
	}
	if unwrapText(wrapText("\001")) != "" {
		t.Fatal(`unwrapText: wrapText("\001")) != ""`)
	}

	if file, _ := unwrapFile([]byte{}); file != "" {
		t.Fatal("unwrapFile([]byte{}) = ok")
	}
	if file, _ := unwrapFile([]byte{cIsFile, 0x01}); file != "" {
		t.Fatal("unwrapFile([]byte{cIsFile, 0x01}) = ok")
	}

	if isText([]byte{}) {
		t.Fatal("isText([]byte{}) = ok")
	}
	if isFile([]byte{}) {
		t.Fatal("isFile([]byte{}) = ok")
	}

	wt := wrapText(tcText)
	if !isText(wt) {
		t.Fatal("wrapText: !isText(wt)")
	}
	if isFile(wt) {
		t.Fatal("wrapText: isFile(wt)")
	}
	if wt[0] != cIsText {
		t.Fatal("wrapText:  wt[0] != cIsText")
	}
	if unwrapText(wt) != tcTextEscaped {
		t.Fatal("wrapText: unwrapText(wt, true) != tcTextEscaped")
	}

	wf := wrapFile(tcFile, []byte(tcText))
	if !isFile(wf) {
		t.Fatal("wrapFile: !isFile(wf)")
	}
	if isText(wf) {
		t.Fatal("wrapFile: isText(wf)")
	}
	if wf[0] != cIsFile || wf[len(tcFile)+1] != cIsFile {
		t.Fatal("wrapFile: wf[0] != cIsFile || wf[len(tcFile)+1] != cIsFile")
	}
	if f, b := unwrapFile(wf); f != tcFile || b != base64.StdEncoding.EncodeToString([]byte(tcText)) {
		t.Fatal("wrapText: f, b := unwrapFile(wf, false); f != tcFile || b != base64.StdEncoding.EncodeToString([]byte(tcText))")
	}

	wfx := wrapFile(tcFile+"<b>some</b>"+".txt", []byte(tcText))
	if f, b := unwrapFile(wfx); f != tcFileEscaped || b != base64.StdEncoding.EncodeToString([]byte(tcText)) {
		t.Fatal("wrapText: f, b := unwrapFile(wf, true); f != tcFileEscaped || b != base64.StdEncoding.EncodeToString([]byte(tcText))")
	}

	if f, b := unwrapFile([]byte{1}); f != "" || b != "" {
		t.Fatal(`wrapFile: f, b := unwrapFile([]byte{1}, false); f != "" || b != ""`)
	}
	wf2 := wrapFile(tcFile+"\x01"+".txt", []byte(tcText))
	if f, b := unwrapFile(wf2); f != "" || b != "" {
		t.Fatal(`wrapFile: f, b := unwrapFile(wf2, false); f != "" || b != ""`)
	}
	wf3 := wrapFile(tcFile, []byte{})
	if f, b := unwrapFile(wf3); f != "" || b != "" {
		t.Fatal(`wrapFile: f, b := unwrapFile(wf3, false); f != "" || b != ""`)
	}
}
