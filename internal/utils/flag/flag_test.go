package flag

import "testing"

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
	_ = GetStringFlagValue(argsSlice, []string{"key"}, "_")
}

func TestBoolFlagValue(t *testing.T) {
	t.Parallel()

	argsSlice := []string{
		"--key", "123",
		"-name",
		"value", "571",
	}

	if !GetBoolFlagValue(argsSlice, []string{"key"}) {
		t.Error("!key")
		return
	}

	if !GetBoolFlagValue(argsSlice, []string{"123"}) {
		t.Error("!123")
		return
	}

	if !GetBoolFlagValue(argsSlice, []string{"name"}) {
		t.Error("!name")
		return
	}

	if !GetBoolFlagValue(argsSlice, []string{"value"}) {
		t.Error("!value")
		return
	}

	if !GetBoolFlagValue(argsSlice, []string{"571"}) {
		t.Error("!571")
		return
	}

	if GetBoolFlagValue(argsSlice, []string{"undefined"}) {
		t.Error("success get undefined value")
		return
	}
}

func TestFlagValue(t *testing.T) {
	t.Parallel()

	argsSlice := []string{
		"--key", "123",
		"-name", "number",
		"-null", "some-value",
		"value", "571",
		"asdfg=12345",
		"-qwerty=67890",
		"--zxcvb=!@#$%",
	}

	if GetStringFlagValue(argsSlice, []string{"key"}, "1") != "123" {
		t.Error("key != 123")
		return
	}

	if GetStringFlagValue(argsSlice, []string{"name"}, "2") != "number" {
		t.Error("name != number")
		return
	}

	if GetStringFlagValue(argsSlice, []string{"value"}, "3") != "571" {
		t.Error("value != 571")
		return
	}

	if GetStringFlagValue(argsSlice, []string{"asdfg"}, "4") != "12345" {
		t.Error("asdfg != 12345")
		return
	}

	if GetStringFlagValue(argsSlice, []string{"qwerty"}, "5") != "67890" {
		t.Error("qwerty != 67890")
		return
	}

	if GetStringFlagValue(argsSlice, []string{"zxcvb"}, "6") != "!@#$%" {
		t.Error("zxcvb != !@#$%")
		return
	}

	if GetStringFlagValue(argsSlice, []string{"unknown"}, "7") != "7" {
		t.Error("unknown != 7")
		return
	}
}
