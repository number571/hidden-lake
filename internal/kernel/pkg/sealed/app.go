package sealed

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2/hybrid"
	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/build"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/network/adapters"
	"github.com/number571/hidden-lake/pkg/network/adapters/tcp"
)

const (
	cAPIListenAddr = "127.0.0.1:0"
	cReadChunkSize = 4096
)

var (
	_ adapters.IRunnerAdapter = &sRelayAdapter{}
)

// sRelayAdapter repeats the behaviour of the standard hlk + hla-tcp pair
// (hla-tcp runExtRelayer, database_enabled=false by default): every message
// received from TCP is broadcast back to all TCP connections. Without it a relay
// could distinguish hlk-sealed within one queue period. Loops are stopped by
// the hash cache of the go-peer network node.
type sRelayAdapter struct {
	fAdapter adapters.IRunnerAdapter
}

func newRelayAdapter(pAdapter adapters.IRunnerAdapter) adapters.IRunnerAdapter {
	return &sRelayAdapter{fAdapter: pAdapter}
}

func (p *sRelayAdapter) Run(pCtx context.Context) error {
	return p.fAdapter.Run(pCtx)
}

func (p *sRelayAdapter) Produce(pCtx context.Context, pMsg layer1.IMessage) error {
	return p.fAdapter.Produce(pCtx, pMsg)
}

func (p *sRelayAdapter) Consume(pCtx context.Context) (layer1.IMessage, error) {
	msg, err := p.fAdapter.Consume(pCtx)
	if err != nil {
		return nil, err
	}
	// Relay synchronously, exactly like hla-tcp runExtRelayer: a stalled send
	// blocks the consume loop (back-pressure) and never drops a message, instead
	// of the earlier goroutine-per-message that a flood turned unbounded.
	_ = p.fAdapter.Produce(pCtx, msg)
	return msg, nil
}

// Run executes hlk-sealed:
//  1. reads exactly one line (boot JSON) from pStdin; errors are wrapped by ErrBoot;
//  2. starts the HLK node (TCP adapter without listen, relays from the config),
//     the seen DB and the internal API on 127.0.0.1:0;
//  3. writes "PRIVKEY k" (only for a generated key), "PUBKEY k", "READY addr"
//     to pStdout; nothing is written to pStdout after READY, logs go to pStderr;
//  4. stops on EOF of pStdin or cancel of pCtx (returns nil) or on a node/API
//     failure (returns error).
func Run(pCtx context.Context, pStdin io.Reader, pStdout, pStderr io.Writer) error {
	line, err := readBootLine(pStdin)
	if err != nil {
		return errors.Join(ErrBoot, ErrReadBoot)
	}
	boot, err := ParseBoot(line)
	clear(line) // secrets: wipe the raw line
	if err != nil {
		return errors.Join(ErrBoot, err)
	}

	seenDB, err := NewSeenDB(boot.GetSeenDBPath(), boot.GetSeenKey())
	if err != nil {
		return err
	}
	defer func() { _ = seenDB.Close() }()

	logging := boot.GetConfig().GetLogging()
	anonLogger := newLogger(pStderr, logging, anon_logger.GetLogFunc())
	httpLogger := newLogger(pStderr, logging, http_logger.GetLogFunc())

	node, tcpAdapter, err := newNode(boot, seenDB, anonLogger)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", cAPIListenAddr)
	if err != nil {
		return errors.Join(ErrListenAPI, err)
	}
	addr := listener.Addr().String()

	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	server := NewAPIServer(NewAPIHandler(
		ctx,
		boot.GetConfig(),
		node,
		tcpAdapter.GetConnKeeper().GetNetworkNode(),
		boot.GetAPIToken(),
		addr,
		httpLogger,
	))

	chErr := make(chan error, 2)
	chNodeDone := make(chan struct{})
	go func() {
		defer close(chNodeDone)
		if err := node.Run(ctx); err != nil && ctx.Err() == nil {
			chErr <- errors.Join(ErrRunNode, err)
		}
	}()
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			chErr <- errors.Join(ErrServeAPI, err)
		}
	}()
	defer func() {
		_ = server.Close()
		cancel()
		<-chNodeDone // node does not use seen DB after this
	}()

	// listener accepts connections since net.Listen: READY is true here
	if err := writeHandshake(pStdout, boot, addr); err != nil {
		return errors.Join(ErrWriteStdout, err)
	}

	// stdin watcher: EOF (parent closed the pipe or died) = stop
	go func() {
		_, _ = io.Copy(io.Discard, pStdin)
		cancel()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-chErr:
		return err
	}
}

// readBootLine reads pR up to and including the first '\n' without a long-lived
// buffer: every buffer that held bytes of the line is wiped. Bytes after '\n'
// in the last read chunk are dropped (only EOF matters after the boot line).
func readBootLine(pR io.Reader) ([]byte, error) {
	chunk := make([]byte, cReadChunkSize)
	defer clear(chunk)

	line := make([]byte, 0, cReadChunkSize)
	for {
		n, err := pR.Read(chunk)
		i := bytes.IndexByte(chunk[:n], '\n')
		if i >= 0 {
			n = i + 1
		}
		if len(line)+n > cap(line) {
			grown := make([]byte, len(line), 2*(len(line)+n))
			copy(grown, line)
			clear(line)
			line = grown
		}
		line = append(line, chunk[:n]...)
		if i >= 0 {
			return line, nil
		}
		if err != nil {
			clear(line)
			return nil, err
		}
	}
}

func writeHandshake(pStdout io.Writer, pBoot IBoot, pAddr string) error {
	privKey := pBoot.GetPrivKey()
	out := ""
	if pBoot.IsPrivKeyGenerated() {
		out += fmt.Sprintf("PRIVKEY %s\n", privKey.ToString())
	}
	out += fmt.Sprintf("PUBKEY %s\nREADY %s\n", privKey.GetPubKey().ToString(), pAddr)
	_, err := io.WriteString(pStdout, out)
	return err
}

// newLogger honours the configured levels: a disabled level gets a nil writer
// (go-peer treats it as off), so info-level routing metadata is not emitted to
// pOut unless the config enables info.
func newLogger(pOut io.Writer, pLogging std_logger.ILogging, pLogFunc gopeer_logger.ILogFunc) gopeer_logger.ILogger {
	sett := &gopeer_logger.SSettings{}
	if pLogging.HasInfo() {
		sett.FInfo = pOut
	}
	if pLogging.HasWarn() {
		sett.FWarn = pOut
	}
	if pLogging.HasErro() {
		sett.FErro = pOut
	}
	return gopeer_logger.NewLogger(gopeer_logger.NewSettings(sett), pLogFunc)
}

// newNode builds the HLK node as internal/kernel/pkg/app/init_node.go does,
// but with the TCP adapter (outgoing connections to relays only) in process.
func newNode(
	pBoot IBoot,
	pSeenDB database.IKVDatabase,
	pLogger gopeer_logger.ILogger,
) (network.IHiddenLakeNode, tcp.ITCPAdapter, error) {
	var (
		cfg           = pBoot.GetConfig()
		cfgSettings   = cfg.GetSettings()
		buildSettings = build.GetSettings()
	)

	adapterSettings := adapters.NewSettings(&adapters.SSettings{
		FWorkSizeBits:     cfgSettings.GetWorkSizeBits(),
		FNetworkKey:       cfgSettings.GetNetworkKey(),
		FMessageSizeBytes: cfgSettings.GetMessageSizeBytes(),
	})

	scheme, err := hybrid.NewScheme(pBoot.GetPrivKey(), adapterSettings.GetMessageSizeBytes())
	if err != nil {
		return nil, nil, errors.Join(ErrCreateNode, err)
	}

	tcpAdapter := tcp.NewTCPAdapter(
		tcp.NewSettings(&tcp.SSettings{
			FAdapterSettings: adapterSettings,
			FServeSettings: &tcp.SServeSettings{
				FAddress:      "", // no listen: outgoing connections only
				FConnNumLimit: buildSettings.FNetworkManager.FConnNumLimit,
			},
		}),
		cache.NewLRUCache(buildSettings.FStorageManager.FCacheHashesCap),
		cfg.GetEndpoints,
	).WithLogger(hlk_settings.GetAppShortNameFMT(), pLogger)

	node, err := network.NewHiddenLakeNode(
		network.NewSettings(&network.SSettings{
			FAdapterSettings: adapterSettings,
			FQBPSettings: &network.SQBPSettings{
				FQueuePeriod:  cfgSettings.GetQueuePeriod(),
				FFetchTimeout: cfgSettings.GetFetchTimeout(),
				FPowParallel:  cfgSettings.GetPowParallel(),
				FQBPConsumers: cfgSettings.GetQBPConsumers(),
			},
			FServeSettings: &network.SServeSettings{
				FServiceName: hlk_settings.GetAppShortNameFMT(),
				FLogger:      pLogger,
			},
		}),
		scheme,
		layer2.NewKeysContainer(),
		pSeenDB,
		newRelayAdapter(tcpAdapter),
		NewServiceHandler(cfg, pBoot.GetPushToken(), pLogger),
	)
	if err != nil {
		return nil, nil, errors.Join(ErrCreateNode, err)
	}

	keys := node.GetOriginNode().GetKeysContainer()
	for _, k := range cfg.GetFriends() {
		if ok := keys.Add(k); !ok {
			return nil, nil, ErrCreateNode
		}
	}

	return node, tcpAdapter, nil
}
