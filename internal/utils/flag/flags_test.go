package flag

import (
	"testing"
)

var (
	tgFlags = NewFlagsBuilder(
		NewFlagBuilder("v", "version").
			WithDescription("print information about service"),
		NewFlagBuilder("p", "path").
			WithDescription("set path to config, database files").
			WithDefaultValue("."),
	).Build()
)

func TestPanicFlagsBuilder(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	_ = NewFlagsBuilder(
		NewFlagBuilder("v", "version").
			WithDescription("print information about service"),
		NewFlagBuilder("v", "version").
			WithDescription("print information about service"),
	).Build()
}

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
