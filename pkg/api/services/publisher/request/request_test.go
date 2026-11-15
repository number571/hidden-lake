package request

import (
	"testing"
)

func TestRequest(t *testing.T) {
	t.Parallel()

	_ = NewMsgboxPushRequest("", []byte{})
	_ = NewMsgboxLoadRequestByUID("", false)
	_ = NewMsgboxLoadRequestByEnum(0, false)
	_ = NewMsgboxSizeRequest(false)
}
