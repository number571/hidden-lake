// nolint: goerr113
package handler

import (
	"testing"

	"github.com/number571/hidden-lake/internal/applications/messenger/internal/msgbroker"
)

func TestFriendsChatWS(t *testing.T) {
	t.Parallel()

	msgBroker := msgbroker.NewMessageBroker()
	handler := FriendsChatWS(msgBroker)

	_ = handler
}
