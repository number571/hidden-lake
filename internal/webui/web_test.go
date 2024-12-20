package webui

import (
	"testing"
)

func TestPanicParseTemplate(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = MustParseTemplate("__unknown__.html")
}

func TestPath(_ *testing.T) {
	_ = MustGetStaticPath()
	_ = MustParseTemplate("index.html")
}
