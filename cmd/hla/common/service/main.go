package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/storage/cache"
)

func main() {
	if len(os.Args) != 3 {
		panic("./service [addr] [logger]")
	}

	addr := os.Args[1]
	logOn := (os.Args[2] == "true")

	iter := uint64(0)
	cache := cache.NewLRUCache(16)

	http.HandleFunc("/last", lastPage(logOn, cache, &iter))
	http.HandleFunc("/load", loadPage(logOn, cache))
	http.HandleFunc("/push", pushPage(logOn, cache, &iter))

	_ = http.ListenAndServe(addr, nil) //nolint:gosec
}

func lastPage(logOn bool, cache cache.ILRUCache, iter *uint64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if logOn {
			log.Printf("PATH: %s; METHOD: %s;\n", r.URL.Path, r.Method)
		}

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		_, _ = fmt.Fprintf(w, "%d.%d", iter, cache.GetIndex())
	}
}

func loadPage(logOn bool, cache cache.ILRUCache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if logOn {
			log.Printf("PATH: %s; METHOD: %s;\n", r.URL.Path, r.Method)
		}

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

func pushPage(logOn bool, cache cache.ILRUCache, iter *uint64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if logOn {
			log.Printf("PATH: %s; METHOD: %s;\n", r.URL.Path, r.Method)
		}

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
