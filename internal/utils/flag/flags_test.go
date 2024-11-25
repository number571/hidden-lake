package flag

import (
	"testing"
)

var (
	tgFlags = NewFlags(
		NewFlagBuilder("v", "version").
			WithDescription("print information about service").
			Build(),
		NewFlagBuilder("p", "path").
			WithDescription("set path to config, database files").
			WithDefaultValue(".").
			Build(),
	)
)

func TestFlagsValidate(t *testing.T) {
	t.Parallel()

	if ok := tgFlags.Validate([]string{"p"}); ok {
		t.Error("success with void string value")
		return
	}
	if ok := tgFlags.Validate([]string{"q"}); ok {
		t.Error("success with not found flag")
		return
	}
	if ok := tgFlags.Validate([]string{"p", "."}); !ok {
		t.Error("failed with success string flag")
		return
	}
	if ok := tgFlags.Validate([]string{"v"}); !ok {
		t.Error("failed with success bool flag")
		return
	}
}
