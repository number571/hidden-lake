package adapted

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/storage/cache"
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

func TestAdaptedConsumer(t *testing.T) {
	t.Parallel()

	iter := uint64(0)
	cache := cache.NewLRUCache(16)

	mux := http.NewServeMux()
	mux.HandleFunc("/last", lastPage(cache, &iter))
	mux.HandleFunc("/load", loadPage(cache))
	mux.HandleFunc("/push", pushPage(cache, &iter))

	addr := testutils.TgAddrs[48]
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

	consumer := NewAdaptedConsumer(sett, "http://"+addr)
	_, err := consumer.Consume(context.Background())
	if err != nil && err.Error() != "status code: 404" {
		t.Error(err)
		return
	}
}

func lastPage(cache cache.ILRUCache, iter *uint64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		_, _ = fmt.Fprintf(w, "%d.%d", iter, cache.GetIndex())
	}
}

func loadPage(cache cache.ILRUCache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()
		strID := query.Get("id")

		id, err := strconv.ParseUint(strID, 10, 64)
		if err != nil {
			fmt.Fprint(w, "!decode data_id")
			return
		}

		key, ok := cache.GetKey(id)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		val, ok := cache.Get(key)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_, _ = w.Write(val)
	}
}

func pushPage(cache cache.ILRUCache, iter *uint64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		res, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		if cache.GetIndex() == 0 {
			*iter++
		}

		cache.Set(hashing.NewHasher(res).ToBytes(), res)
	}
}
