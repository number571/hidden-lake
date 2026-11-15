package sealed

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/pkg/network/request"
)

const (
	tcPushToken = "push-token-0123456789abcdef0123456789abcdef"
	tcPushBody  = `{"text":"hello"}`
)

type tsPush struct {
	method string
	path   string
	body   string
	header http.Header
}

func testLogger() gopeer_logger.ILogger {
	return gopeer_logger.NewLogger(
		gopeer_logger.NewSettings(&gopeer_logger.SSettings{}),
		func(_ gopeer_logger.ILogArg) string { return "" },
	)
}

func testServiceServer(t *testing.T) (*httptest.Server, chan tsPush, *atomic.Int64) {
	t.Helper()

	pushes := make(chan tsPush, 1)
	hits := &atomic.Int64{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		body, _ := io.ReadAll(r.Body)
		pushes <- tsPush{method: r.Method, path: r.URL.Path, body: string(body), header: r.Header.Clone()}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("response must not go to network"))
	}))
	t.Cleanup(srv.Close)
	return srv, pushes, hits
}

func testServiceConfig(t *testing.T, pServiceAddr string) IConfig {
	t.Helper()

	cfg, err := NewConfig("", pServiceAddr, nil, map[string]string{
		"alice": tgPrivKey1.GetPubKey().ToString(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func testServiceRequest(pHost string) request.IRequest {
	return request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(pHost).
		WithPath("/push").
		WithHead(map[string]string{
			"Hlk-Auth":        "peer-forged-token",
			"Hlk-Sender-Name": "forged-alias",
			"Hlk-Other":       "x",
			"X-Custom":        "y",
			"Content-Type":    "text/plain",
		}).
		WithBody([]byte(tcPushBody)).
		Build()
}

func TestServiceHandler(t *testing.T) {
	t.Parallel()

	srv, pushes, _ := testServiceServer(t)
	cfg := testServiceConfig(t, strings.TrimPrefix(srv.URL, "http://"))
	handleF := NewServiceHandler(cfg, tcPushToken, testLogger())

	resp, err := handleF(context.Background(), tgPrivKey1.GetPubKey(), testServiceRequest(tcMessenger))
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		t.Fatal("service response returned to network")
	}

	push := <-pushes
	if push.method != http.MethodPost || push.path != "/push" || push.body != tcPushBody {
		t.Fatalf("invalid push: %+v", push)
	}
	if push.header.Get("Hlk-Auth") != tcPushToken {
		t.Fatal("invalid Hlk-Auth")
	}
	if v := push.header.Values("Hlk-Sender-Name"); len(v) != 1 || v[0] != "alice" {
		t.Fatal("invalid Hlk-Sender-Name")
	}
	for _, h := range []string{"Hlk-Other", "X-Custom", "Content-Type"} {
		if push.header.Get(h) != "" {
			t.Fatalf("peer header %s passed", h)
		}
	}
	if len(push.header.Values("Hlk-Auth")) != 1 {
		t.Fatal("peer Hlk-Auth passed")
	}
}

func TestServiceHandlerReject(t *testing.T) {
	t.Parallel()

	srv, _, hits := testServiceServer(t)
	cfg := testServiceConfig(t, strings.TrimPrefix(srv.URL, "http://"))
	handleF := NewServiceHandler(cfg, tcPushToken, testLogger())

	_, err := handleF(context.Background(), tgPrivKey1.GetPubKey(), testServiceRequest("hls-filesharer"))
	if !errors.Is(err, ErrUndefinedService) {
		t.Fatal("unknown host accepted")
	}

	_, err = handleF(context.Background(), tgPrivKey2.GetPubKey(), testServiceRequest(tcMessenger))
	if !errors.Is(err, ErrUnknownSender) {
		t.Fatal("unknown sender accepted")
	}

	if hits.Load() != 0 {
		t.Fatal("rejected request reached service")
	}
}

func TestServiceHandlerUnavailable(t *testing.T) {
	t.Parallel()

	srv, _, _ := testServiceServer(t)
	addr := strings.TrimPrefix(srv.URL, "http://")
	srv.Close()

	handleF := NewServiceHandler(testServiceConfig(t, addr), tcPushToken, testLogger())
	_, err := handleF(context.Background(), tgPrivKey1.GetPubKey(), testServiceRequest(tcMessenger))
	if !errors.Is(err, ErrRequestToService) {
		t.Fatal("unavailable service is not an error")
	}
}

func TestServiceHandlerInvalidRequest(t *testing.T) {
	t.Parallel()

	srv, _, hits := testServiceServer(t)

	handleF := NewServiceHandler(testServiceConfig(t, strings.TrimPrefix(srv.URL, "http://")), tcPushToken, testLogger())

	// Only an exact POST /push is forwarded; anything else is rejected and the
	// service is never contacted.
	testCases := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "push"},
		{http.MethodPost, "/other"},
		{http.MethodPost, "/push/more"},
		{http.MethodPost, "/push?a"},
		{http.MethodPost, ""},
		{http.MethodGet, "/push"},
		{http.MethodPut, "/push"},
	}
	for _, tc := range testCases {
		req := request.NewRequestBuilder().
			WithMethod(tc.method).
			WithHost(tcMessenger).
			WithPath(tc.path).
			WithBody([]byte(tcPushBody)).
			Build()
		if _, err := handleF(context.Background(), tgPrivKey1.GetPubKey(), req); !errors.Is(err, ErrInvalidRequest) {
			t.Fatalf("%s %q accepted", tc.method, tc.path)
		}
	}

	if hits.Load() != 0 {
		t.Fatal("invalid request reached the service")
	}
}

func TestServiceHandlerNoRedirect(t *testing.T) {
	t.Parallel()

	other, _, otherHits := testServiceServer(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL+"/push", http.StatusTemporaryRedirect)
	}))
	t.Cleanup(srv.Close)

	handleF := NewServiceHandler(testServiceConfig(t, strings.TrimPrefix(srv.URL, "http://")), tcPushToken, testLogger())
	if _, err := handleF(context.Background(), tgPrivKey1.GetPubKey(), testServiceRequest(tcMessenger)); err != nil {
		t.Fatal(err)
	}
	if otherHits.Load() != 0 {
		t.Fatal("redirect followed")
	}
}
