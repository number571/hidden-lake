package sealed

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/adapters/tcp"
)

const (
	tcFastNetworkKey = "test-sealed-fast"
	tcFastMsgSize    = 8192
	tcWaitTimeout    = 60 * time.Second
)

func TestMain(m *testing.M) {
	// hybrid test network with a short queue period (not the default network)
	networks := build.GetNetworks()
	networks[tcFastNetworkKey] = build.SNetwork{
		FCryptoSchemeType: "hybrid",
		FMessageSizeBytes: tcFastMsgSize,
		FWorkSizeBits:     0,
		FFetchTimeoutMS:   5_000,
		FQueuePeriodMS:    200,
	}
	if err := build.SetNetworks(networks); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestRunBootErrors(t *testing.T) {
	t.Parallel()

	secret := strings.Repeat("5e", 32)
	testCases := []struct {
		name  string
		stdin string
	}{
		{"empty stdin", ""},
		{"no newline", `{"api_token":"` + secret + `"}`},
		{"invalid json with secret", `{"api_token":"` + secret + `",` + "\n"},
		{"invalid boot with secret", `{"api_token":"` + secret + `","push_token":"x"}` + "\n"},
	}

	for _, tc := range testCases {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		err := Run(context.Background(), strings.NewReader(tc.stdin), stdout, stderr)
		if !errors.Is(err, ErrBoot) {
			t.Fatalf("%s: expected boot error, got %v", tc.name, err)
		}
		if stdout.Len() != 0 {
			t.Fatalf("%s: stdout is not empty", tc.name)
		}
		if strings.Contains(err.Error(), secret) || strings.Contains(stderr.String(), secret) {
			t.Fatalf("%s: secret in error", tc.name)
		}
	}
}

func TestRunSeenDBError(t *testing.T) {
	t.Parallel()

	m := testBootMap()
	m["network_key"] = tcFastNetworkKey
	m["seen_db_path"] = filepath.Join(t.TempDir(), "missing", "seen.db")

	stdout := &bytes.Buffer{}
	err := Run(context.Background(), bytes.NewReader(testBootLine(t, m)), stdout, io.Discard)
	if err == nil || errors.Is(err, ErrBoot) || !errors.Is(err, ErrOpenSeenDB) {
		t.Fatalf("expected seen db error, got %v", err)
	}
	if stdout.Len() != 0 {
		t.Fatal("stdout is not empty")
	}
}

type tsNode struct {
	stdin   *io.PipeWriter
	done    chan error
	extra   chan int64
	privKey string
	pubKey  string
	addr    string
	token   string
}

func testStartNode(t *testing.T, pBoot map[string]any) *tsNode {
	t.Helper()

	stdinR, stdinW := io.Pipe()
	stdoutR, stdoutW := io.Pipe()
	node := &tsNode{
		stdin: stdinW,
		done:  make(chan error, 1),
		extra: make(chan int64, 1),
		token: pBoot["api_token"].(string),
	}

	go func() {
		node.done <- Run(context.Background(), stdinR, stdoutW, io.Discard)
		_ = stdoutW.Close()
	}()
	go func() {
		_, _ = stdinW.Write(testBootLine(t, pBoot))
	}()

	scanner := bufio.NewScanner(stdoutR)
	scanner.Buffer(make([]byte, 0, 64<<10), 64<<10)
	for scanner.Scan() {
		cmd, arg, _ := strings.Cut(scanner.Text(), " ")
		switch cmd {
		case "PRIVKEY":
			node.privKey = arg
		case "PUBKEY":
			node.pubKey = arg
		case "READY":
			node.addr = arg
		default:
			t.Fatalf("unexpected stdout line: %q", scanner.Text())
		}
		if node.addr != "" {
			break
		}
	}
	if node.addr == "" {
		t.Fatalf("no READY: %v", <-node.done)
	}
	go func() { // nothing expected after READY
		n, _ := io.Copy(io.Discard, stdoutR)
		node.extra <- n
	}()

	t.Cleanup(func() { _ = stdinW.Close() })
	return node
}

func (p *tsNode) request(t *testing.T, pMethod, pPath, pBody string) (int, string) {
	t.Helper()

	req, err := http.NewRequestWithContext(context.Background(), pMethod, "http://"+p.addr+pPath, strings.NewReader(pBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(CHeaderAuth, p.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

func (p *tsNode) waitOnline(t *testing.T) {
	t.Helper()

	deadline := time.Now().Add(tcWaitTimeout)
	for time.Now().Before(deadline) {
		if code, body := p.request(t, http.MethodGet, "/api/network/online", ""); code == http.StatusOK && body != "[]" && body != "[]\n" {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatal("node is not connected to relay")
}

func testFreeAddr(t *testing.T) string {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	return addr
}

// testStartRelay runs a TCP relay: every received message is broadcast to all connections.
func testStartRelay(t *testing.T) string {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	addr := testFreeAddr(t)
	relay := tcp.NewTCPAdapter(
		tcp.NewSettings(&tcp.SSettings{
			FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
				FMessageSizeBytes: tcFastMsgSize,
				FNetworkKey:       tcFastNetworkKey,
			}),
			FServeSettings: &tcp.SServeSettings{FAddress: addr},
		}),
		cache.NewLRUCache(1024),
		func() []string { return nil },
	)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() { defer wg.Done(); _ = relay.Run(ctx) }()
	go func() {
		defer wg.Done()
		for {
			msg, err := relay.Consume(ctx)
			if err != nil {
				return
			}
			_ = relay.Produce(ctx, msg)
		}
	}()
	t.Cleanup(func() { cancel(); wg.Wait() })

	deadline := time.Now().Add(tcWaitTimeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			_ = conn.Close()
			return addr
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("relay is not listening")
	return ""
}

type tsPushRecord struct {
	body   string
	sender string
	auth   string
}

func TestRunIntegration(t *testing.T) {
	t.Parallel()

	relay := testStartRelay(t)

	pushes := make(chan tsPushRecord, 4)
	service := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		pushes <- tsPushRecord{body: string(body), sender: r.Header.Get("Hlk-Sender-Name"), auth: r.Header.Get(CHeaderAuth)}
		w.Header().Set("Hlk-Response-Mode", "off")
	}))
	t.Cleanup(service.Close)

	privB := asymmetric.NewPrivKey()
	pushB := strings.Repeat("b7", 32)

	bootA := testBootMap()
	bootA["priv_key"] = ""
	bootA["network_key"] = tcFastNetworkKey
	bootA["friends"] = map[string]string{"bob": privB.GetPubKey().ToString()}
	bootA["endpoints"] = []string{relay}
	bootA["service_addr"] = "127.0.0.1:1"
	bootA["api_token"] = strings.Repeat("a1", 32)
	bootA["seen_db_path"] = filepath.Join(t.TempDir(), "seen.db")

	bootB := testBootMap()
	bootB["priv_key"] = privB.ToString()
	bootB["network_key"] = tcFastNetworkKey
	bootB["endpoints"] = []string{relay}
	bootB["service_addr"] = strings.TrimPrefix(service.URL, "http://")
	bootB["api_token"] = strings.Repeat("c2", 32)
	bootB["push_token"] = pushB
	bootB["seen_db_path"] = filepath.Join(t.TempDir(), "seen.db")

	nodeA := testStartNode(t, bootA)
	if nodeA.privKey == "" || asymmetric.LoadPrivKey(nodeA.privKey) == nil {
		t.Fatal("generated PRIVKEY is not printed")
	}
	if nodeA.pubKey != asymmetric.LoadPrivKey(nodeA.privKey).GetPubKey().ToString() {
		t.Fatal("PUBKEY does not match PRIVKEY")
	}
	bootB["friends"] = map[string]string{"alice": nodeA.pubKey}

	nodeB := testStartNode(t, bootB)
	if nodeB.privKey != "" || nodeB.pubKey != privB.GetPubKey().ToString() {
		t.Fatal("invalid handshake of B")
	}

	// READY address is the address of the API (Host check passes)
	if code, body := nodeA.request(t, http.MethodGet, "/api/config/settings", ""); code != http.StatusOK || !strings.Contains(body, "payload_size_bytes") {
		t.Fatalf("settings by READY address: %d %s", code, body)
	}

	nodeA.waitOnline(t)
	nodeB.waitOnline(t)

	reqBody := `{"method":"POST","host":"hls-messenger","path":"/push","body":"` +
		base64.StdEncoding.EncodeToString([]byte("hi")) + `"}`
	if code, body := nodeA.request(t, http.MethodPost, "/api/network/request?friend=bob", reqBody); code != http.StatusOK {
		t.Fatalf("send request: %d %s", code, body)
	}

	select {
	case push := <-pushes:
		if push.body != "hi" || push.sender != "alice" || push.auth != pushB {
			t.Fatalf("invalid push: %+v", push)
		}
	case <-time.After(tcWaitTimeout):
		t.Fatal("push not received")
	}

	// EOF of stdin stops the node quickly
	select {
	case err := <-nodeA.done:
		t.Fatalf("run stopped before EOF: %v", err)
	case <-time.After(500 * time.Millisecond):
	}
	if code, _ := nodeA.request(t, http.MethodGet, "/api/network/online", ""); code != http.StatusOK {
		t.Fatal("api is not available before EOF")
	}
	start := time.Now()
	_ = nodeA.stdin.Close()
	select {
	case err := <-nodeA.done:
		if err != nil {
			t.Fatalf("run after EOF: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("run does not stop after EOF")
	}
	t.Logf("stopped after EOF in %s", time.Since(start))
	if n := <-nodeA.extra; n != 0 {
		t.Fatalf("%d bytes written to stdout after READY", n)
	}

	_ = nodeB.stdin.Close()
	select {
	case <-nodeB.done:
	case <-time.After(5 * time.Second):
		t.Fatal("run B does not stop after EOF")
	}
}

// testStartObserver runs a listening TCP adapter that does not relay anything.
func testStartObserver(t *testing.T) (string, tcp.ITCPAdapter) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	addr := testFreeAddr(t)
	observer := tcp.NewTCPAdapter(
		tcp.NewSettings(&tcp.SSettings{
			FAdapterSettings: adapters.NewSettings(&adapters.SSettings{
				FMessageSizeBytes: tcFastMsgSize,
				FNetworkKey:       tcFastNetworkKey,
			}),
			FServeSettings: &tcp.SServeSettings{FAddress: addr},
		}),
		cache.NewLRUCache(1024),
		func() []string { return nil },
	)
	done := make(chan struct{})
	go func() { defer close(done); _ = observer.Run(ctx) }()
	t.Cleanup(func() { cancel(); <-done })

	deadline := time.Now().Add(tcWaitTimeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err == nil {
			_ = conn.Close()
			return addr, observer
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("observer is not listening")
	return "", nil
}

// TestRunRelaysMessages: like hla-tcp, a message received from one connection
// is broadcast to all connections (otherwise the node is distinguishable).
func TestRunRelaysMessages(t *testing.T) {
	t.Parallel()

	addr1, observer1 := testStartObserver(t)
	addr2, observer2 := testStartObserver(t)

	boot := testBootMap()
	boot["network_key"] = tcFastNetworkKey
	boot["endpoints"] = []string{addr1, addr2}
	boot["seen_db_path"] = filepath.Join(t.TempDir(), "seen.db")
	node := testStartNode(t, boot)

	deadline := time.Now().Add(tcWaitTimeout)
	for {
		_, body := node.request(t, http.MethodGet, "/api/network/online", "")
		if strings.Contains(body, addr1) && strings.Contains(body, addr2) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("node is not connected to observers")
		}
		time.Sleep(100 * time.Millisecond)
	}

	// foreign message (not decryptable by the node)
	msg := layer1.NewMessage(
		layer1.NewConstructSettings(&layer1.SConstructSettings{
			FSettings: layer1.NewSettings(&layer1.SSettings{FNetworkKey: tcFastNetworkKey}),
		}),
		random.NewRandom().GetBytes(tcFastMsgSize),
	)
	if err := observer1.Produce(context.Background(), msg); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), tcWaitTimeout)
	defer cancel()
	for {
		got, err := observer2.Consume(ctx)
		if err != nil {
			t.Fatal("foreign message is not relayed by the node")
		}
		if bytes.Equal(got.GetHash(), msg.GetHash()) {
			return
		}
	}
}

type tsChunkReader struct {
	data  []byte
	chunk int
}

func (p *tsChunkReader) Read(pB []byte) (int, error) {
	if len(p.data) == 0 {
		return 0, io.EOF
	}
	n := min(len(pB), p.chunk, len(p.data))
	copy(pB, p.data[:n])
	p.data = p.data[n:]
	return n, nil
}

func TestReadBootLine(t *testing.T) {
	t.Parallel()

	long := strings.Repeat("x", 3*cReadChunkSize+17)
	testCases := []struct {
		name  string
		input string
		chunk int
		line  string
		rest  string
	}{
		{"short", "abc\nrest", 1 << 20, "abc\n", "rest"},
		{"long in chunks", long + "\nrest", 1000, long + "\n", "rest"},
		{"byte by byte", "abc\nrest", 1, "abc\n", "rest"},
		{"newline only", "\n", 10, "\n", ""},
	}
	for _, tc := range testCases {
		r := &tsChunkReader{data: []byte(tc.input), chunk: tc.chunk}
		line, err := readBootLine(r)
		if err != nil || string(line) != tc.line {
			t.Fatalf("%s: invalid line (%v)", tc.name, err)
		}
		// bytes after the line are not consumed beyond the current chunk
		rest, _ := io.ReadAll(r)
		if !strings.HasSuffix(tc.rest, string(rest)) {
			t.Fatalf("%s: invalid rest %q", tc.name, rest)
		}
	}

	for _, input := range []string{"", "no newline"} {
		if _, err := readBootLine(strings.NewReader(input)); err == nil {
			t.Fatalf("line without newline accepted: %q", input)
		}
	}
}

// tsRelayFake is an IRunnerAdapter whose Produce blocks, to prove the relay
// serializes rebroadcasts through one worker instead of a goroutine per message.
type tsRelayFake struct {
	fProduced int64
	fBlock    chan struct{}
}

func (p *tsRelayFake) Run(context.Context) error { return nil }

func (p *tsRelayFake) Consume(pCtx context.Context) (layer1.IMessage, error) {
	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	default:
		return nil, nil
	}
}
func (p *tsRelayFake) Produce(pCtx context.Context, _ layer1.IMessage) error {
	atomic.AddInt64(&p.fProduced, 1)
	select {
	case <-p.fBlock:
	case <-pCtx.Done():
	}
	return nil
}

func TestRelayAdapterSynchronousBackpressure(t *testing.T) {
	t.Parallel()

	fake := &tsRelayFake{fBlock: make(chan struct{})}
	relay := newRelayAdapter(fake)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, err := relay.Consume(ctx); err != nil {
				return
			}
		}
	}()
	time.Sleep(100 * time.Millisecond)

	// A blocked relay Produce must block the consume loop (back-pressure), so
	// exactly one Produce is in flight. Before the fix each Consume spawned its
	// own goroutine and the count ran away.
	if n := atomic.LoadInt64(&fake.fProduced); n != 1 {
		t.Fatalf("relay ran %d concurrent Produce calls; expected 1 (synchronous back-pressure)", n)
	}
	cancel()
	<-done
}

func TestNewLoggerHonoursLevels(t *testing.T) {
	t.Parallel()

	logging, err := std_logger.LoadLogging([]string{std_logger.CLogWarn, std_logger.CLogErro})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	lg := newLogger(&buf, logging, func(pArg gopeer_logger.ILogArg) string { return pArg.(string) })

	lg.PushInfo("info-line")
	if buf.Len() != 0 {
		t.Fatalf("info written though disabled: %q", buf.String())
	}
	lg.PushErro("erro-line")
	if !strings.Contains(buf.String(), "erro-line") {
		t.Fatal("erro level not written")
	}
}
