package std

import "testing"

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestLogging(t *testing.T) {
	t.Parallel()

	logging, err := LoadLogging([]string{"info", "erro"})
	if err != nil {
		t.Fatal(err)
	}
	if !logging.HasInfo() {
		t.Fatal("failed has info")
	}
	if logging.HasWarn() {
		t.Fatal("failed has warn")
	}
	if !logging.HasErro() {
		t.Fatal("failed has erro")
	}
	if _, err := LoadLogging([]string{"info", "unknown"}); err == nil {
		t.Fatal("success load invalid logging")
	}
}
