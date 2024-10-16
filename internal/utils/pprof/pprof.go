package pprof

import (
	"net/http"
	"time"

	"net/http/pprof"
)

func InitPprofService(pAddr string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	server := &http.Server{
		Addr:        pAddr,
		Handler:     http.TimeoutHandler(mux, time.Minute, "timeout"),
		ReadTimeout: (5 * time.Second),
	}

	return server
}
