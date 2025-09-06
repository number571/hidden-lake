// nolint: err113
package closer

import (
	"errors"
	"io"
	"testing"
)

type tsCloser struct {
	fFlag bool
}

func TestCloser(t *testing.T) {
	t.Parallel()

	closer := NewCloser(
		testNewCloser(false),
		testNewCloser(false),
		testNewCloser(false),
	)

	if err := closer.Close(); err != nil {
		t.Fatal(err)
	}

	if err := NewCloser(testNewCloser(true)).Close(); err == nil {
		t.Fatal("nothing error?")
	}
}

func testNewCloser(flag bool) io.Closer {
	return &tsCloser{flag}
}

func (c *tsCloser) Close() error {
	if c.fFlag {
		return errors.New("some error")
	}
	return nil
}
