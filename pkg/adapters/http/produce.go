package http

import (
	"io"
	"net/http"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/build"
)

func (p *sHTTPAdapter) produceHandler() func(http.ResponseWriter, *http.Request) {
	adapterSettings := p.fSettings.GetAdapterSettings()
	cache := cache.NewLRUCache(build.GSettings.FNetworkManager.FCacheHashesCap)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		msgSize := adapterSettings.GetMessageSizeBytes()
		msgLen := uint64(msgSize+net_message.CMessageHeadSize) << 1 // nolint: unconvert
		msgStr := make([]byte, msgLen)
		n, err := io.ReadFull(r.Body, msgStr)
		if err != nil || uint64(n) != msgLen {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		msg, err := net_message.LoadMessage(adapterSettings, string(msgStr))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hash := msg.GetHash()
		if ok := cache.Set(hash, []byte{}); !ok {
			w.WriteHeader(http.StatusAlreadyReported)
			return
		}

		p.fNetMsgChan <- msg
	}
}
