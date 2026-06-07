package anon

import (
	"context"
	"net"
	"testing"
	"time"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/conn"
)

var (
	_ conn.IConn = &tsConn{}
	_ net.Conn   = &tsNetConn{}
	_ net.Addr   = &tsAddr{}
)

type tsConn struct{}

func (p *tsConn) Close() error                                        { return nil }
func (p *tsConn) GetSettings() conn.ISettings                         { return nil }
func (p *tsConn) GetSocket() net.Conn                                 { return &tsNetConn{} }
func (p *tsConn) WriteMessage(context.Context, layer1.IMessage) error { return nil }
func (p *tsConn) ReadMessage(context.Context, chan<- struct{}) (layer1.IMessage, error) {
	return nil, nil
}

type tsNetConn struct{}

func (p *tsNetConn) Read(_ []byte) (n int, err error)   { return 0, nil }
func (p *tsNetConn) Write(_ []byte) (n int, err error)  { return 0, nil }
func (p *tsNetConn) Close() error                       { return nil }
func (p *tsNetConn) LocalAddr() net.Addr                { return &tsAddr{} }
func (p *tsNetConn) RemoteAddr() net.Addr               { return &tsAddr{} }
func (p *tsNetConn) SetDeadline(_ time.Time) error      { return nil }
func (p *tsNetConn) SetReadDeadline(_ time.Time) error  { return nil }
func (p *tsNetConn) SetWriteDeadline(_ time.Time) error { return nil }

type tsAddr struct{}

func (p *tsAddr) Network() string { return "tcp" }
func (p *tsAddr) String() string  { return "192.168.0.1:2000" }

const (
	tcService = "TST"
	tcHash    = "96cb1f0968adba001ebc216708a02c8d2817b1a77fad1206012c22716a9b130b"
	tcFmtLog  = "service=TST type=ENQRQ hash=96CB1F09...00000000 proof=12345 size=1024B conn=127.0.0.1"
)

func TestLoggerPanic(t *testing.T) {
	t.Parallel()

	logFunc := GetLogFunc()
	for i := 0; i < 3; i++ {
		testLoggerPanic(t, logFunc, i)
	}
}

func testLoggerPanic(t *testing.T, f logger.ILogFunc, n int) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()
	switch n {
	case 0:
		f(struct{}{})
	case 1:
		logger := testNewAnonLogger()
		f(logger) // without type
	case 2:
		logger := testNewAnonLogger().WithType(255)
		f(logger) // with unknown type
	}
}

func TestLogger(t *testing.T) {
	t.Parallel()

	logger := testNewAnonLogger().
		WithType(anon_logger.CLogBaseEnqueueRequest)

	logFunc := GetLogFunc()
	if l := logFunc(logger); l != tcFmtLog {
		t.Log(l)
		t.Fatal("result fmtLog != tcFmtLog")
	}
}

func testNewAnonLogger() anon_logger.ILogBuilder {
	return anon_logger.NewLogBuilder(tcService).
		WithHash(encoding.HexDecode(tcHash)).
		WithProof(12345).
		WithSize(1024).
		WithConn("127.0.0.1")
}
