package adapted

import (
	"context"
	"net/http"
	"testing"
	"time"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAdaptedError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestAdaptedProducer(t *testing.T) {
	t.Parallel()

	statusOK := true

	mux := http.NewServeMux()
	mux.HandleFunc("/push", func(w http.ResponseWriter, _ *http.Request) {
		if statusOK {
			return
		}
		w.WriteHeader(500)
	})

	addr := testutils.TgAddrs[47]
	srv := &http.Server{
		Addr:        addr,
		Handler:     mux,
		ReadTimeout: time.Second,
	}
	defer srv.Close()
	go func() { _ = srv.ListenAndServe() }()

	time.Sleep(200 * time.Millisecond)

	sett := net_message.NewSettings(&net_message.SSettings{
		FWorkSizeBits: 1,
		FNetworkKey:   "_",
	})

	msg := net_message.NewMessage(
		net_message.NewConstructSettings(&net_message.SConstructSettings{
			FSettings: sett,
			FParallel: 1,
		}),
		payload.NewPayload32(1, []byte("hello, world!")),
	)

	producer := NewAdaptedProducer("http://" + addr)
	if err := producer.Produce(context.Background(), msg); err != nil {
		t.Error(err)
		return
	}

	statusOK = false
	if err := producer.Produce(context.Background(), msg); err == nil {
		t.Error("success produce with status code not ok")
		return
	}

	if err := NewAdaptedProducer("").Produce(context.Background(), msg); err == nil {
		t.Error("success produce with invalid addr")
		return
	}
}
