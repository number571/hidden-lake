package flag

import (
	"testing"
)

var (
	tgFlags = NewFlagsBuilder(
		NewFlagBuilder("-v", "--version").
			WithDescription("print information about application"),
		NewFlagBuilder("-p", "--path").
			WithDescription("set path to config, database files").
			WithDefinedValue("."),
	).Build()
)

func TestPanicFlagsGet(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	_ = NewFlagsBuilder(
		NewFlagBuilder("-v", "--version").
			WithDescription("print information about application"),
	).Build().Get("--unknown")
}

func TestPanicFlagsBuilder(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	_ = NewFlagsBuilder(
		NewFlagBuilder("-v", "--version").
			WithDescription("print information about application"),
		NewFlagBuilder("-v", "--version").
			WithDescription("print information about application"),
	).Build()
}

func TestFlagsValidate(t *testing.T) {
	t.Parallel()

	if ok := tgFlags.Validate([]string{"-p"}); ok {
		t.Fatal("success with void string value")
	}
	if ok := tgFlags.Validate([]string{"-q"}); ok {
		t.Fatal("success with not found flag")
	}
	if ok := tgFlags.Validate([]string{"-p", "."}); !ok {
		t.Fatal("failed with success string flag")
	}
	if ok := tgFlags.Validate([]string{"-v"}); !ok {
		t.Fatal("failed with success bool flag")
	}
}
