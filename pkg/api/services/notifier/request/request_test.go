package request

import (
	"testing"
)

func TestRequest(t *testing.T) {
	t.Parallel()

	_ = NewNotifyRequest([]byte{})
}
