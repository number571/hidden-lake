package flag

import (
	"testing"
)

func TestPanicFlagValue(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()
	argsSlice := []string{
		"--key",
	}
	_ = NewFlagBuilder("--key").Build().GetStringValue(argsSlice)
}

func TestPanicInt64Value(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()
	argsSlice := []string{
		"--key", "qwerty",
		"--key2",
	}
	_ = NewFlagBuilder("--key").Build().GetInt64Value(argsSlice)
	_ = NewFlagBuilder("--key2").Build().GetInt64Value(argsSlice)
}

func TestInt64Value(t *testing.T) {
	t.Parallel()

	argsSlice := []string{
		"--key", "123",
		"-name", "number",
		"-null", "some-value",
		"value", "571",
	}

	if NewFlagBuilder("--key").Build().GetInt64Value(argsSlice) != 123 {
		t.Fatal("key != 123")
	}
	if NewFlagBuilder("value").Build().GetInt64Value(argsSlice) != 571 {
		t.Fatal("value != 571")
	}
}

func TestBoolFlagValue(t *testing.T) {
	t.Parallel()

	argsSlice := []string{
		"--key", "123",
		"-name",
		"value", "571",
	}

	if !NewFlagBuilder("--key").Build().GetBoolValue(argsSlice) {
		t.Fatal("!key")
	}

	if !NewFlagBuilder("123").Build().GetBoolValue(argsSlice) {
		t.Fatal("!123")
	}

	if !NewFlagBuilder("-name").Build().GetBoolValue(argsSlice) {
		t.Fatal("!name")
	}

	if !NewFlagBuilder("value").Build().GetBoolValue(argsSlice) {
		t.Fatal("!value")
	}

	if !NewFlagBuilder("571").Build().GetBoolValue(argsSlice) {
		t.Fatal("!571")
	}

	if NewFlagBuilder("undefined").Build().GetBoolValue(argsSlice) {
		t.Fatal("success get undefined value")
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
		t.Fatal("key != 123")
	}

	if NewFlagBuilder("-name").Build().GetStringValue(argsSlice) != "number" {
		t.Fatal("name != number")
	}

	if NewFlagBuilder("value").Build().GetStringValue(argsSlice) != "571" {
		t.Fatal("value != 571")
	}

	if NewFlagBuilder("unknown").WithDefinedValue("7").Build().GetStringValue(argsSlice) != "7" {
		t.Fatal("unknown != 7")
	}
}
