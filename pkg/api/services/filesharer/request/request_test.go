package request

import "testing"

func TestRequest(t *testing.T) {
	t.Parallel()

	_ = NewInfoRequest("", false)
	_ = NewListRequest(0, false)
	_ = NewLoadRequest("", 0, false)
}
