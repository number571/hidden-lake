package std

import (
	"testing"
)

var (
	_ = &tsLogger{}
)

type tsLogger struct{}

func (l *tsLogger) HasInfo() bool {
	return true
}

func (l *tsLogger) HasWarn() bool {
	return true
}

func (l *tsLogger) HasErro() bool {
	return true
}

func TestGetLogFunc(t *testing.T) {
	t.Parallel()

	f := GetLogFunc()
	if l := f("string"); l != "string" {
		t.Fatal("incorrect logger work")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()
	_ = f(struct{}{})
}
