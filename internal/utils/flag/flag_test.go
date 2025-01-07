package flag

import (
	"testing"
)

func TestPanicFlagValue(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()
	argsSlice := []string{
		"--key",
	}
	_ = NewFlagBuilder("--key").Build().GetStringValue(argsSlice)
}

func TestBoolFlagValue(t *testing.T) {
	t.Parallel()

	argsSlice := []string{
		"--key", "123",
		"-name",
		"value", "571",
	}

	if !NewFlagBuilder("--key").Build().GetBoolValue(argsSlice) {
		t.Error("!key")
		return
	}

	if !NewFlagBuilder("123").Build().GetBoolValue(argsSlice) {
		t.Error("!123")
		return
	}

	if !NewFlagBuilder("-name").Build().GetBoolValue(argsSlice) {
		t.Error("!name")
		return
	}

	if !NewFlagBuilder("value").Build().GetBoolValue(argsSlice) {
		t.Error("!value")
		return
	}

	if !NewFlagBuilder("571").Build().GetBoolValue(argsSlice) {
		t.Error("!571")
		return
	}

	if NewFlagBuilder("undefined").Build().GetBoolValue(argsSlice) {
		t.Error("success get undefined value")
		return
	}
}

func TestStringFlagValue(t *testing.T) {
	t.Parallel()

	argsSlice := []string{
		"--key", "123",
		"-name", "number",
		"-null", "some-value",
		"value", "571",
	}

	if NewFlagBuilder("--key").Build().GetStringValue(argsSlice) != "123" {
		t.Error("key != 123")
		return
	}

	if NewFlagBuilder("-name").Build().GetStringValue(argsSlice) != "number" {
		t.Error("name != number")
		return
	}

	if NewFlagBuilder("value").Build().GetStringValue(argsSlice) != "571" {
		t.Error("value != 571")
		return
	}

	if NewFlagBuilder("unknown").WithDefinedValue("7").Build().GetStringValue(argsSlice) != "7" {
		t.Error("unknown != 7")
		return
	}
}
