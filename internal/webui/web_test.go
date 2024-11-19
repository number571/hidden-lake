package webui

import (
	"testing"
)

func TestPath(_ *testing.T) {
	_ = GetStaticPath()
	_ = GetTemplatePath()
	_ = MustParseTemplate("index.html")
}
