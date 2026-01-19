// nolint: err113
package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	testutils "github.com/number571/hidden-lake/test/utils"
)

type tsRequest tsResponse
type tsResponse struct {
	FMessage string `json:"message"`
}

type tsReadCloser struct {
	io.ReadCloser
}

func (p *tsReadCloser) Read(_ []byte) (n int, err error) { return 0, errors.New("some error") }
func (p *tsReadCloser) Close() error                     { return nil }

type tsResponseWriter struct{}

func (p *tsResponseWriter) Header() http.Header         { return make(http.Header) }
func (p *tsResponseWriter) Write(_ []byte) (int, error) { return 0, errors.New("some error") }
func (p *tsResponseWriter) WriteHeader(_ int)           {}

const (
	tcMessage          = "hello, world!"
	tcMethodNotAllowed = "method not allowed"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestResponseAPI(t *testing.T) {
	t.Parallel()

	if err := Response(&tsResponseWriter{}, 200, []byte{123}); err == nil {
		t.Fatal("success response with invalid response writer")
	}
}

func TestLoadResponse(t *testing.T) {
	t.Parallel()

	if _, err := loadResponse(0, &tsReadCloser{}); err == nil {
		t.Fatal("success load response with invalid readCloser")
	}
}

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[0]
	unknownURL := "http://" + addr + "/unknown"

	client := &http.Client{
		Timeout: time.Minute / 4,
	}

	srv := testRunServer(addr)
	defer func() { _ = srv.Close() }()

	if _, err := Request(context.Background(), client, http.MethodGet, addr, nil); err == nil {
		t.Fatal("success request on incorrect url address")
	}

	if _, err := Request(context.Background(), client, http.MethodGet, unknownURL, nil); err == nil {
		t.Fatal("success request on unknown url address")
	}
}

func TestRequestResponseAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[1]
	testURL := "http://" + addr + "/test"

	client := &http.Client{
		Timeout: time.Minute / 4,
	}

	srv := testRunServer(addr)
	defer func() { _ = srv.Close() }()

	if _, err := Request(context.Background(), client, http.MethodGet, "\n\t\a", nil); err == nil {
		t.Fatal("success request on invalid url")
	}

	if _, err := Request(context.Background(), client, http.MethodPatch, testURL, nil); err == nil {
		t.Fatal("PATCH: success request on method not allowed")
	}

	respGET, err := Request(context.Background(), client, http.MethodGet, testURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	x := new(tsResponse)
	if err := encoding.DeserializeJSON(respGET, x); err != nil {
		t.Fatal(err)
	}

	if x.FMessage != tcMessage {
		t.Fatal("GET: got message is invalid")
	}

	// bytes
	respPOST1, err := Request(context.Background(), client, http.MethodPost, testURL, []byte(tcMessage))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(respPOST1, bytes.Join([][]byte{[]byte("echo"), []byte(tcMessage)}, []byte{1})) {
		t.Fatal("POST1: got message is invalid")
	}

	// string
	respPOST2, err := Request(context.Background(), client, http.MethodPost, testURL, tcMessage)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(respPOST2, bytes.Join([][]byte{[]byte("echo"), []byte(tcMessage)}, []byte{1})) {
		t.Fatal("POST2: got message is invalid")
	}

	// struct
	respPOST3, err := Request(context.Background(), client, http.MethodPost, testURL, tsRequest{FMessage: tcMessage})
	if err != nil {
		t.Fatal(err)
	}

	msg := fmt.Sprintf(`{"message":"%s"}`, tcMessage)
	if !bytes.Equal(respPOST3, bytes.Join([][]byte{[]byte("echo"), []byte(msg)}, []byte{1})) {
		t.Fatal("POST3: got message is invalid")
	}
}

func testRunServer(addr string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			_ = Response(w, http.StatusMethodNotAllowed, tcMethodNotAllowed)
			return
		}

		if r.Method == http.MethodGet {
			_ = Response(w, http.StatusOK, tsResponse{FMessage: tcMessage})
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		_ = Response(w, http.StatusOK, bytes.Join([][]byte{[]byte("echo"), data}, []byte{1}))
	})

	srv := &http.Server{
		Addr:        addr,
		Handler:     http.TimeoutHandler(mux, time.Minute/2, "timeout"),
		ReadTimeout: time.Second,
	}

	go func() { _ = srv.ListenAndServe() }()
	time.Sleep(200 * time.Millisecond)
	return srv
}
