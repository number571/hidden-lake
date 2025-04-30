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

	err := CloseAll([]io.Closer{
		testNewCloser(false),
		testNewCloser(false),
		testNewCloser(false),
	})
	if err != nil {
		t.Error(err)
		return
	}

	if err := CloseAll([]io.Closer{testNewCloser(true)}); err == nil {
		t.Error("nothing error?")
		return
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
