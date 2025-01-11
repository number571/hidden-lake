package tcp

import "testing"

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestSettings(t *testing.T) {
	t.Parallel()

	_ = NewSettings(nil)
}

func TestTCPAdapter(t *testing.T) {
	t.Parallel()

}
