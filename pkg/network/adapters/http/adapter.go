package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/broker"
	internal_anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"

	"github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
	"github.com/number571/hidden-lake/pkg/api/adapters/http/client"
)

var (
	_ IHTTPAdapter = &sHTTPAdapter{}
)

type sHTTPAdapter struct {
	fSettings   ISettings
	fNetMsgChan chan layer1.IMessage

	fConnsGetter func() []string
	fOnlines     *sOnlines
	fCache       cache.ICache

	fShortName  string
	fLogger     logger.ILogger
	fDataBroker broker.IDataBroker
	fHandlers   map[string]http.HandlerFunc
}

type sOnlines struct {
	fMutex sync.RWMutex
	fSlice []string
}

func NewHTTPAdapter(
	pSettings ISettings,
	pCache cache.ICache,
	pConnsGetter func() []string,
) IHTTPAdapter {
	return &sHTTPAdapter{
		fSettings:    pSettings,
		fCache:       pCache,
		fNetMsgChan:  make(chan layer1.IMessage, pSettings.GetChannelSize()),
		fConnsGetter: pConnsGetter,
		fOnlines:     &sOnlines{fSlice: pConnsGetter()},
		fLogger: logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
		fDataBroker: broker.NewDataBroker(
			pSettings.GetChannelSize(),
			pSettings.GetConnNumLimit(),
		),
	}
}

func (p *sHTTPAdapter) WithHandlers(pHandlers map[string]http.HandlerFunc) IHTTPAdapter {
	p.fHandlers = pHandlers
	return p
}

func (p *sHTTPAdapter) WithLogger(pName string, pLogger logger.ILogger) IHTTPAdapter {
	p.fShortName = pName
	p.fLogger = pLogger
	return p
}

func (p *sHTTPAdapter) Run(pCtx context.Context) error {
	go func() {
		_ = p.runSubscriber(pCtx)
		// internal logger
	}()

	address := p.fSettings.GetAddress()
	if address == "" {
		<-pCtx.Done()
		return pCtx.Err()
	}

	mux := http.NewServeMux()

	mux.HandleFunc(settings.CHandleAdapterProducePath, p.adapterProduceHandler(pCtx))
	mux.HandleFunc(settings.CHandleAdapterConsumePath, p.adapterConsumeHandler(pCtx))

	for k, v := range p.fHandlers {
		mux.HandleFunc(k, v)
	}

	httpServer := &http.Server{
		Addr:         address,
		Handler:      http.TimeoutHandler(mux, p.fSettings.GetHandleTimeout(), "handle timeout"),
		ReadTimeout:  p.fSettings.GetReadTimeout(),
		WriteTimeout: p.fSettings.GetHandleTimeout(),
	}
	go func() {
		<-pCtx.Done()
		_ = httpServer.Close()
	}()

	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return errors.Join(ErrRunning, err)
	}
	return context.Canceled
}

func (p *sHTTPAdapter) Produce(pCtx context.Context, pNetMsg layer1.IMessage) error {
	logBuilder := anon_logger.NewLogBuilder(p.fShortName)
	logBuilder.
		WithType(internal_anon_logger.CLogBaseSendNetworkMessage).
		WithHash(pNetMsg.GetHash()).
		WithProof(pNetMsg.GetProof()).
		WithSize(len(pNetMsg.ToBytes())).
		WithConn("http")

	// adapter can redirect received message
	hash := encoding.HexEncode(pNetMsg.GetHash())
	_ = p.fCache.Set(hash, []byte{})
	p.fDataBroker.Produce(pNetMsg)

	connects := p.fConnsGetter()
	if len(connects) == 0 {
		if p.fDataBroker.CountSubscribers() > 0 {
			// produces messages with `adapterConsumeHandler` function
			p.fLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoHasOnlySubscribers))
			return nil
		}
		p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnNoConnections))
		return ErrNoConnections
	}

	N := len(connects)
	errs := make([]error, N)

	wg := &sync.WaitGroup{}
	wg.Add(N)
	for i, url := range connects {
		go func(i int, url string) {
			defer wg.Done()
			errs[i] = p.produceMessage(pCtx, url, pNetMsg)
		}(i, url)
	}
	wg.Wait()

	onlines := make([]string, 0, N)
	for i := range errs {
		if errs[i] != nil {
			continue
		}
		onlines = append(onlines, connects[i])
	}

	p.fOnlines.fMutex.Lock()
	p.fOnlines.fSlice = onlines
	p.fOnlines.fMutex.Unlock()

	if joinedErr := errors.Join(errs...); joinedErr != nil {
		p.fLogger.PushWarn(logBuilder)
		return joinedErr
	}

	p.fLogger.PushInfo(logBuilder)
	return nil
}

func (p *sHTTPAdapter) Consume(pCtx context.Context) (layer1.IMessage, error) {
	for {
		select {
		case <-pCtx.Done():
			return nil, pCtx.Err()
		case msg := <-p.fNetMsgChan:
			return msg, nil
		}
	}
}

func (p *sHTTPAdapter) runSubscriber(pCtx context.Context) error {
	connListener := func(addr string, closed chan struct{}) {
		logBuilder := anon_logger.NewLogBuilder(p.fShortName)
		logBuilder.WithConn(addr)

		for {
			select {
			case <-pCtx.Done():
				return
			case <-closed:
				return
			default:
				msg, err := p.consumeMessage(pCtx, addr)
				if err != nil {
					p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogBaseRecvNetworkMessage))
					select {
					case <-pCtx.Done():
					case <-time.After(time.Second):
					}
					continue
				}

				logBuilder.
					WithHash(msg.GetHash()).
					WithProof(msg.GetProof()).
					WithSize(len(msg.ToBytes()))

				p.fLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseRecvNetworkMessage))

				hash := encoding.HexEncode(msg.GetHash())
				if ok := p.fCache.Set(hash, []byte{}); !ok {
					continue
				}

				if ok := p.pushMessageToChan(msg); !ok {
					p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnMessageChanOverflow))
					continue
				}
			}
		}
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	mapConns := make(map[string]chan struct{})
	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		case <-ticker.C:
			conns := p.fConnsGetter()
			mapCheck := make(map[string]struct{}, len(conns))

			// enrich map
			for _, c := range conns {
				mapCheck[c] = struct{}{}
				if _, ok := mapConns[c]; ok {
					continue
				}
				ch := make(chan struct{}, 1)
				mapConns[c] = ch
				go connListener(c, ch)
			}

			// clear map
			for k, v := range mapConns {
				if _, ok := mapCheck[k]; !ok {
					close(v)
					delete(mapConns, k)
				}
			}
		}
	}
}

func (p *sHTTPAdapter) produceMessage(pCtx context.Context, pHost string, pMsg layer1.IMessage) error {
	hlaClient := client.NewClient(
		client.NewRequester(
			pHost,
			&http.Client{Timeout: p.fSettings.GetHandleTimeout()},
			p.fSettings.GetAdapterSettings(),
		),
	)
	return hlaClient.ProduceMessage(pCtx, pMsg)
}

func (p *sHTTPAdapter) consumeMessage(pCtx context.Context, pHost string) (layer1.IMessage, error) {
	hlaClient := client.NewClient(
		client.NewRequester(
			pHost,
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
			p.fSettings.GetAdapterSettings(),
		),
	)
	return hlaClient.ConsumeMessage(pCtx, p.fSettings.GetSubscribeID())
}

func (p *sHTTPAdapter) GetOnlines() []string {
	p.fOnlines.fMutex.RLock()
	defer p.fOnlines.fMutex.RUnlock()

	return p.fOnlines.fSlice
}

func (p *sHTTPAdapter) adapterProduceHandler(_ context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		adapterSettings := p.fSettings.GetAdapterSettings()

		logBuilder := anon_logger.NewLogBuilder(p.fShortName)
		logBuilder.WithConn(r.RemoteAddr)

		if r.Method != http.MethodPost {
			p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnInvalidRequestMethod))
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		msgLen := adapterSettings.GetMessageSizeBytes() + layer1.CMessageHeadSize
		msgLen <<= 1 // message hex_encoded
		msgStr := make([]byte, msgLen)
		n, err := io.ReadFull(r.Body, msgStr)
		if err != nil || uint64(n) != msgLen { //nolint:gosec
			p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnFailedReadFullBytes))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		msg, err := layer1.LoadMessage(adapterSettings, string(msgStr))
		if err != nil {
			p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnMessageNull))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logBuilder.
			WithHash(msg.GetHash()).
			WithProof(msg.GetProof()).
			WithSize(len(msg.ToBytes()))

		if msg.GetPayload().GetHead() != build.GetSettings().FProtoMask.FNetwork {
			p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnPayloadNull))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p.fLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseRecvNetworkMessage))

		hash := encoding.HexEncode(msg.GetHash())
		if ok := p.fCache.Set(hash, []byte{}); !ok {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		if ok := p.pushMessageToChan(msg); !ok {
			p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnMessageChanOverflow))
		}
	}
}

func (p *sHTTPAdapter) pushMessageToChan(pMsg layer1.IMessage) bool {
	p.fDataBroker.Produce(pMsg)
	select {
	case p.fNetMsgChan <- pMsg:
		return true
	default:
		return false
	}
}

func (p *sHTTPAdapter) adapterConsumeHandler(pCtx context.Context) func(w http.ResponseWriter, r *http.Request) {
	buildSettings := build.GetSettings()

	return func(w http.ResponseWriter, r *http.Request) {
		logBuilder := anon_logger.NewLogBuilder(p.fShortName)
		logBuilder.WithConn(r.RemoteAddr)

		if r.Method != http.MethodGet {
			p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnInvalidRequestMethod))
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		sid := r.URL.Query().Get("sid")

		if err := p.fDataBroker.Register(sid); err != nil {
			p.fLogger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnLimitOfSubscribers))
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		ctx, cancel := context.WithTimeout(pCtx, buildSettings.GetHttpReadTimeout())
		defer cancel()

		v, err := p.fDataBroker.Consume(ctx, sid)
		if err != nil {
			p.fLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoNoContent))
			_ = api.Response(w, http.StatusNoContent, []byte{})
			return
		}

		msg, ok := v.(layer1.IMessage)
		if !ok {
			p.fLogger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroInvalidMessageType))
			_ = api.Response(w, http.StatusInternalServerError, []byte{})
			return
		}

		logBuilder.
			WithHash(msg.GetHash()).
			WithProof(msg.GetProof()).
			WithSize(len(msg.ToBytes()))

		p.fLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseSendNetworkMessage))
		_ = api.Response(w, http.StatusOK, msg.ToString())
	}
}
