package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/number571/hidden-lake/build"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/cache"
)

const (
	netMessageChanSize = 32
)

var (
	_ IHTTPAdapter = &sHTTPAdapter{}
)

type sHTTPAdapter struct {
	fNetMsgChan  chan net_message.IMessage
	fHTTPServer  *http.Server
	fConnsGetter func() []string
}

func NewHTTPAdapter(pSettings ISettings, pConnsGetter func() []string) IHTTPAdapter {
	msgChan := make(chan net_message.IMessage, netMessageChanSize)
	cache := cache.NewLRUCache(build.GSettings.FNetworkManager.FCacheHashesCap)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		msgLen := uint64(pSettings.GetMessageSizeBytes()+net_message.CMessageHeadSize) << 1 // nolint: unconvert
		msgStr := make([]byte, msgLen)
		n, err := io.ReadFull(r.Body, msgStr)
		if err != nil || uint64(n) != msgLen {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		msg, err := net_message.LoadMessage(pSettings, string(msgStr))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hash := msg.GetHash()
		if _, ok := cache.Get(hash); ok {
			w.WriteHeader(http.StatusAlreadyReported)
			return
		}
		_ = cache.Set(hash, []byte{})
		msgChan <- msg
	})

	return &sHTTPAdapter{
		fNetMsgChan: msgChan,
		fHTTPServer: &http.Server{
			Addr:        pSettings.GetAddress(),
			Handler:     http.TimeoutHandler(mux, 5*time.Second, "timeout"),
			ReadTimeout: (5 * time.Second),
		},
		fConnsGetter: pConnsGetter,
	}
}

func (p *sHTTPAdapter) Run(pCtx context.Context) error {
	go func() {
		<-pCtx.Done()
		p.fHTTPServer.Close()
	}()
	return p.fHTTPServer.ListenAndServe()
}

func (p *sHTTPAdapter) Produce(pCtx context.Context, pNetMsg net_message.IMessage) error {
	connects := p.fConnsGetter()
	N := len(connects)
	errs := make([]error, N)

	wg := &sync.WaitGroup{}
	wg.Add(N)
	for i, url := range connects {
		go func(i int, url string) {
			defer wg.Done()
			req, err := http.NewRequestWithContext(
				pCtx,
				http.MethodPost,
				"http://"+url,
				strings.NewReader(pNetMsg.ToString()),
			)
			if err != nil {
				errs[i] = err
				return
			}
			client := &http.Client{Timeout: 5 * time.Second}
			rsp, err := client.Do(req)
			if err != nil {
				errs[i] = err
				return
			}
			rsp.Body.Close()
		}(i, url)
	}
	wg.Wait()

	return errors.Join(errs...)
}

func (p *sHTTPAdapter) Consume(pCtx context.Context) (net_message.IMessage, error) {
	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	case msg := <-p.fNetMsgChan:
		return msg, nil
	}
}
