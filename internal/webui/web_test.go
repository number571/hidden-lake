package webui

import (
	"testing"
)

func TestPath(_ *testing.T) {
	_ = GetTemplatePath()
	_ = GetStaticPath()
}
