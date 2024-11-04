package adapted

import "testing"

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAdaptedError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestAdaptedProducer(t *testing.T) {
	t.Parallel()

	_ = NewAdaptedProducer("")

}
