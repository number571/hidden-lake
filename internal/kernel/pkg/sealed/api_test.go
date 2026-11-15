package sealed

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/internal/utils/database"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/adapters/tcp"
)

const (
	tcAPIAddr  = "127.0.0.1:47001"
	tcAPIToken = "api-token-0123456789abcdef0123456789abcdef"
)

type tsAPI struct {
	handler http.Handler
	cfg     IConfig
	node    network.IHiddenLakeNode
}

type tsReadSpy struct {
	reads atomic.Int64
}

func (p *tsReadSpy) Read(_ []byte) (int, error) {
	p.reads.Add(1)
	return 0, io.EOF
}

func testNewAPI(t *testing.T) *tsAPI {
	t.Helper()

	cfg, err := NewConfig("", "127.0.0.1:1", nil, map[string]string{
		"alice": tgPrivKey2.GetPubKey().ToString(),
	})
	if err != nil {
		t.Fatal(err)
	}

	adapterSettings := adapters.NewSettings(&adapters.SSettings{
		FMessageSizeBytes: cfg.GetSettings().GetMessageSizeBytes(),
	})
	scheme, err := hybrid.NewScheme(tgPrivKey1, adapterSettings.GetMessageSizeBytes())
	if err != nil {
		t.Fatal(err)
	}

	tcpAdapter := tcp.NewTCPAdapter(
		tcp.NewSettings(&tcp.SSettings{FAdapterSettings: adapterSettings}),
		cache.NewLRUCache(1024),
		cfg.GetEndpoints,
	)
	node, err := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{FAdapterSettings: adapterSettings}),
		scheme,
		layer2.NewKeysContainer(),
		database.NewVoidKVDatabase(),
		tcpAdapter,
		NewServiceHandler(cfg, tcPushToken, testLogger()),
	)
	if err != nil {
		t.Fatal(err)
	}

	handler := NewAPIHandler(
		context.Background(),
		cfg,
		node,
		tcpAdapter.GetConnKeeper().GetNetworkNode(),
		tcAPIToken,
		tcAPIAddr,
		testLogger(),
	)
	return &tsAPI{handler: handler, cfg: cfg, node: node}
}

func (p *tsAPI) do(t *testing.T, pMethod, pPath string, pBody io.Reader) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(pMethod, "http://"+tcAPIAddr+pPath, pBody)
	req.Header.Set(CHeaderAuth, tcAPIToken)
	rec := httptest.NewRecorder()
	p.handler.ServeHTTP(rec, req)
	return rec
}

func TestAPIUnauthorized(t *testing.T) {
	t.Parallel()

	api := testNewAPI(t)

	testCases := []struct {
		name  string
		host  string
		token string
	}{
		{"no token", tcAPIAddr, ""},
		{"invalid token", tcAPIAddr, tcAPIToken + "x"},
		{"prefix token", tcAPIAddr, tcAPIToken[:len(tcAPIToken)-1]},
		{"foreign host", "foreign.example:47001", tcAPIToken},
		{"localhost host", "localhost:47001", tcAPIToken},
	}

	for _, tc := range testCases {
		spy := &tsReadSpy{}
		req := httptest.NewRequest(http.MethodPost, "http://"+tcAPIAddr+"/api/config/friends", spy)
		req.Host = tc.host
		if tc.token != "" {
			req.Header.Set(CHeaderAuth, tc.token)
		}
		rec := httptest.NewRecorder()
		api.handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("%s: expected 401, got %d", tc.name, rec.Code)
		}
		if spy.reads.Load() != 0 {
			t.Fatalf("%s: body was read", tc.name)
		}
	}

	// body bigger than limit is still 401 without auth
	big := strings.NewReader(strings.Repeat("a", 2*cMaxBodySize))
	req := httptest.NewRequest(http.MethodPost, "http://"+tcAPIAddr+"/api/network/request?friend=alice", big)
	rec := httptest.NewRecorder()
	api.handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatal("big body without auth is not 401")
	}
}

func TestAPIEmptyToken(t *testing.T) {
	t.Parallel()

	api := testNewAPI(t)
	handler := NewAPIHandler(context.Background(), api.cfg, api.node, nil, "", tcAPIAddr, testLogger())

	req := httptest.NewRequest(http.MethodGet, "http://"+tcAPIAddr+"/api/config/settings", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("empty api token accepts requests without header: %d", rec.Code)
	}
}

func TestAPIBodyLimit(t *testing.T) {
	t.Parallel()

	api := testNewAPI(t)

	// valid request padded by leading whitespace up to the limit (+1 byte)
	body := `{"alias_name":"bob","friend_key":"` + tgPrivKey1.GetPubKey().ToString() + `"}`
	padding := cMaxBodySize - len(body)

	rec := api.do(t, http.MethodPost, "/api/config/friends", strings.NewReader(strings.Repeat(" ", padding+1)+body))
	if rec.Code == http.StatusOK {
		t.Fatal("body bigger than limit accepted")
	}
	if _, ok := api.cfg.GetFriends()["bob"]; ok {
		t.Fatal("friend added from oversized body")
	}

	rec = api.do(t, http.MethodPost, "/api/config/friends", strings.NewReader(strings.Repeat(" ", padding)+body))
	if rec.Code != http.StatusOK {
		t.Fatalf("body equal to limit rejected: %d", rec.Code)
	}
}

func TestAPISettings(t *testing.T) {
	t.Parallel()

	api := testNewAPI(t)

	rec := api.do(t, http.MethodGet, "/api/config/settings", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("settings: %d", rec.Code)
	}
	var v map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if size, ok := v["payload_size_bytes"].(float64); !ok || size <= 0 {
		t.Fatal("no payload_size_bytes")
	}
}

func TestAPIFriends(t *testing.T) {
	t.Parallel()

	api := testNewAPI(t)
	keys := api.node.GetOriginNode().GetKeysContainer()
	bobKey := tgPrivKey1.GetPubKey()

	body := `{"alias_name":"bob","friend_key":"` + bobKey.ToString() + `"}`
	if rec := api.do(t, http.MethodPost, "/api/config/friends", strings.NewReader(body)); rec.Code != http.StatusOK {
		t.Fatalf("add friend: %d", rec.Code)
	}
	if _, ok := api.cfg.GetFriends()["bob"]; !ok {
		t.Fatal("friend not added to config")
	}
	if _, ok := keys.Get(bobKey.GetHasher().ToString()); !ok {
		t.Fatal("friend not added to keys container")
	}

	body = `{"alias_name":"bob"}`
	if rec := api.do(t, http.MethodDelete, "/api/config/friends", strings.NewReader(body)); rec.Code != http.StatusOK {
		t.Fatalf("del friend: %d", rec.Code)
	}
	if _, ok := api.cfg.GetFriends()["bob"]; ok {
		t.Fatal("friend not deleted from config")
	}
	if _, ok := keys.Get(bobKey.GetHasher().ToString()); ok {
		t.Fatal("friend not deleted from keys container")
	}
}

func TestAPIConnectsOnline(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			defer func() { _ = c.Close() }()
		}
	}()

	api := testNewAPI(t)
	relay := ln.Addr().String()

	if rec := api.do(t, http.MethodPost, "/api/config/connects", strings.NewReader(relay)); rec.Code != http.StatusOK {
		t.Fatalf("add connect: %d", rec.Code)
	}
	for _, c := range []string{"not-an-addr", ":9581", "127.0.0.1:", "127.0.0.1:0", "127.0.0.1:65536", "127.0.0.1:http", "127.0.0.1:-1"} {
		if rec := api.do(t, http.MethodPost, "/api/config/connects", strings.NewReader(c)); rec.Code != http.StatusBadRequest {
			t.Fatalf("invalid connect %q: %d", c, rec.Code)
		}
	}

	var list []string
	rec := api.do(t, http.MethodGet, "/api/config/connects", nil)
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil || !slices.Equal(list, []string{relay}) {
		t.Fatalf("connects: %s", rec.Body.String())
	}

	rec = api.do(t, http.MethodGet, "/api/network/online", nil)
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil || !slices.Equal(list, []string{relay}) {
		t.Fatalf("online: %s", rec.Body.String())
	}

	if rec := api.do(t, http.MethodDelete, "/api/config/connects", strings.NewReader(relay)); rec.Code != http.StatusOK {
		t.Fatalf("del connect: %d", rec.Code)
	}
	if len(api.cfg.GetEndpoints()) != 0 {
		t.Fatal("connect not deleted from config")
	}

	rec = api.do(t, http.MethodGet, "/api/network/online", nil)
	if strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("online after delete: %s", rec.Body.String())
	}
}

func TestAPINetworkRequest(t *testing.T) {
	t.Parallel()

	api := testNewAPI(t)
	body := `{"method":"POST","host":"hls-messenger","path":"/push","body":"aGVsbG8="}`

	rec := api.do(t, http.MethodPost, "/api/network/request?friend=unknown", strings.NewReader(body))
	if rec.Code < 400 || rec.Code >= 500 {
		t.Fatalf("unknown friend: %d", rec.Code)
	}

	// POST = send without waiting for response (QB queue)
	rec = api.do(t, http.MethodPost, "/api/network/request?friend=alice", strings.NewReader(body))
	if rec.Code != http.StatusOK {
		t.Fatalf("send request: %d %s", rec.Code, rec.Body.String())
	}

	rec = api.do(t, http.MethodGet, "/api/network/request?friend=alice", nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("get request: %d", rec.Code)
	}
}

func TestAPIServer(t *testing.T) {
	t.Parallel()

	srv := NewAPIServer(http.NotFoundHandler())
	if srv.ReadHeaderTimeout <= 0 || srv.ReadTimeout <= 0 || srv.WriteTimeout <= 0 {
		t.Fatal("server timeouts are not set")
	}
}

func TestAPIFriendsConcurrent(t *testing.T) {
	t.Parallel()

	api := testNewAPI(t)

	const n = 16
	keys := make([]string, n)
	for i := range keys {
		keys[i] = asymmetric.NewPrivKey().GetPubKey().ToString()
	}

	wg := sync.WaitGroup{}
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			body := fmt.Sprintf(`{"alias_name":"f%d","friend_key":"%s"}`, i, keys[i])
			if rec := api.do(t, http.MethodPost, "/api/config/friends", strings.NewReader(body)); rec.Code != http.StatusOK {
				t.Errorf("add friend %d: %d", i, rec.Code)
			}
		}()
	}
	wg.Wait()

	if len(api.cfg.GetFriends()) != n+1 {
		t.Fatalf("lost friend updates: %d", len(api.cfg.GetFriends()))
	}
}

func TestAPIFriendsInvalid(t *testing.T) {
	t.Parallel()

	api := testNewAPI(t)

	ffKey := "PubKey{" + strings.Repeat("FF", asymmetric.CPubKeySize) + "}"
	for _, body := range []string{
		`{"alias_name":"bob","friend_key":"` + ffKey + `"}`,
		`{"alias_name":"bo/b","friend_key":"` + tgPrivKey1.GetPubKey().ToString() + `"}`,
	} {
		if rec := api.do(t, http.MethodPost, "/api/config/friends", strings.NewReader(body)); rec.Code == http.StatusOK {
			t.Fatalf("invalid friend accepted: %d", rec.Code)
		}
	}
	if len(api.cfg.GetFriends()) != 1 {
		t.Fatal("invalid friend added")
	}

	// handler still works after a rejected (panicking) key
	body := `{"alias_name":"bob","friend_key":"` + tgPrivKey1.GetPubKey().ToString() + `"}`
	if rec := api.do(t, http.MethodPost, "/api/config/friends", strings.NewReader(body)); rec.Code != http.StatusOK {
		t.Fatalf("add friend after invalid: %d", rec.Code)
	}
}
