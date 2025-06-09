package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/logger"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/go-peer/pkg/storage/cache"
	"github.com/number571/hidden-lake/build"
	internal_anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"
	"github.com/number571/hidden-lake/internal/utils/name"
	hla_client "github.com/number571/hidden-lake/pkg/adapters/http/client"
	"github.com/number571/hidden-lake/pkg/adapters/http/settings"
)

const (
	netMessageChanSize = 32
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

	fShortName string
	fLogger    logger.ILogger
	fHandlers  map[string]http.HandlerFunc
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
		fNetMsgChan:  make(chan layer1.IMessage, netMessageChanSize),
		fConnsGetter: pConnsGetter,
		fOnlines:     &sOnlines{fSlice: pConnsGetter()},
		fLogger: logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
	}
}

func (p *sHTTPAdapter) WithHandlers(pHandlers map[string]http.HandlerFunc) IHTTPAdapter {
	p.fHandlers = pHandlers
	return p
}

func (p *sHTTPAdapter) WithLogger(pName name.IServiceName, pLogger logger.ILogger) IHTTPAdapter {
	p.fShortName = pName.Short()
	p.fLogger = pLogger
	return p
}

func (p *sHTTPAdapter) Run(pCtx context.Context) error {
	address := p.fSettings.GetAddress()
	if address == "" {
		<-pCtx.Done()
		return pCtx.Err()
	}

	mux := http.NewServeMux()
	mux.HandleFunc(settings.CHandleNetworkAdapterPath, p.adapterHandler)
	for k, v := range p.fHandlers {
		mux.HandleFunc(k, v)
	}

	httpServer := &http.Server{
		Addr:         address,
		Handler:      http.TimeoutHandler(mux, p.fSettings.GetHandleTimeout(), "handle timeout"),
		ReadTimeout:  p.fSettings.GetReadTimeout(),
		WriteTimeout: p.fSettings.GetWriteTimeout(),
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

	connects := p.fConnsGetter()
	if len(connects) == 0 {
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
			httpClient := &http.Client{Timeout: p.fSettings.GetHandleTimeout()}
			errs[i] = hla_client.NewClient(
				hla_client.NewRequester(url, httpClient),
			).ProduceMessage(pCtx, pNetMsg)
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
	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	case msg := <-p.fNetMsgChan:
		return msg, nil
	}
}

func (p *sHTTPAdapter) GetOnlines() []string {
	p.fOnlines.fMutex.RLock()
	defer p.fOnlines.fMutex.RUnlock()

	return p.fOnlines.fSlice
}

func (p *sHTTPAdapter) adapterHandler(w http.ResponseWriter, r *http.Request) {
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

	if msg.GetPayload().GetHead() != build.GetSettings().FProtoMask.FNetwork {
		p.fLogger.PushWarn(logBuilder.WithType(anon_logger.CLogWarnPayloadNull))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logBuilder.
		WithHash(msg.GetHash()).
		WithProof(msg.GetProof()).
		WithSize(len(msg.ToBytes()))

	if ok := p.fCache.Set(msg.GetHash(), []byte{}); !ok {
		p.fLogger.PushInfo(logBuilder.WithType(anon_logger.CLogInfoExist))
		w.WriteHeader(http.StatusAccepted)
		return
	}

	p.fLogger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoRecvNetworkMessage))
	p.fNetMsgChan <- msg
}
