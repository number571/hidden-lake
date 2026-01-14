package handler

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/message"
	"github.com/number571/hidden-lake/internal/utils/api"
)

func HandleListenMessageAPI(pCtx context.Context, pBroker message.IMessageBroker) http.HandlerFunc {
	buildSettings := build.GetSettings()

	return func(pW http.ResponseWriter, pR *http.Request) {
		if pR.Method != http.MethodGet {
			_ = api.Response(pW, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		friend := pR.URL.Query().Get("friend")
		sid := pR.URL.Query().Get("sid")

		ctx, cancel := context.WithTimeout(pCtx, buildSettings.GetHttpReadTimeout())
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				_ = api.Response(pW, http.StatusNoContent, []byte{})
				return
			case c, ok := <-pBroker.Consume(sid):
				if !ok {
					_ = api.Response(pW, http.StatusNoContent, []byte{})
					return
				}
				if c.GetFriend() != friend {
					continue
				}
				_ = api.Response(pW, http.StatusOK, c.GetMessage().ToString())
				return
			}
		}
	}
}
