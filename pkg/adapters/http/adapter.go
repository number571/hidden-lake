package http

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	hla_client "github.com/number571/hidden-lake/internal/adapters/pkg/client"
	"github.com/number571/hidden-lake/internal/adapters/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/name"
)

const (
	netMessageChanSize = 32
)

var (
	_ IHTTPAdapter = &sHTTPAdapter{}
)

type sHTTPAdapter struct {
	fSettings   ISettings
	fNetMsgChan chan net_message.IMessage

	fConnsGetter func() []string
	fOnlines     *sOnlines

	fShortName string
	fLogger    logger.ILogger
	fHandlers  []IHandler
}

type sOnlines struct {
	fMutex sync.RWMutex
	fSlice []string
}

func NewHTTPAdapter(
	pSettings ISettings,
	pConnsGetter func() []string,
) IHTTPAdapter {
	return &sHTTPAdapter{
		fSettings:    pSettings,
		fNetMsgChan:  make(chan net_message.IMessage, netMessageChanSize),
		fConnsGetter: pConnsGetter,
		fOnlines:     &sOnlines{fSlice: pConnsGetter()},
		fLogger: logger.NewLogger(
			logger.NewSettings(&logger.SSettings{}),
			func(_ logger.ILogArg) string { return "" },
		),
	}
}

func (p *sHTTPAdapter) WithHandlers(pHandlers ...IHandler) IHTTPAdapter {
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
	mux.HandleFunc(settings.CHandleNetworkAdapterPath, p.produceHandler())
	for _, handler := range p.fHandlers {
		mux.HandleFunc(handler.GetPath(), handler.GetFunc())
	}
	httpServer := &http.Server{
		Addr:        address,
		Handler:     mux,
		ReadTimeout: (5 * time.Second),
	}
	go func() {
		<-pCtx.Done()
		httpServer.Close()
	}()
	return httpServer.ListenAndServe()
}

func (p *sHTTPAdapter) Produce(pCtx context.Context, pNetMsg net_message.IMessage) error {

	connects := p.fConnsGetter()
	N := len(connects)
	errs := make([]error, N)

	wg := &sync.WaitGroup{}
	wg.Add(N)
	for i, url := range connects {
		go func(i int, url string) {
			defer wg.Done()
			errs[i] = hla_client.NewClient(
				hla_client.NewRequester(url, &http.Client{Timeout: 5 * time.Second}),
			).ProduceMessage(pCtx, pNetMsg)
		}(i, url)
	}
	wg.Wait()

	onlines := make([]string, 0, N)
	for i := range errs {
		if errs[i] == nil {
			onlines = append(onlines, connects[i])
		}
	}

	p.fOnlines.fMutex.Lock()
	p.fOnlines.fSlice = onlines
	p.fOnlines.fMutex.Unlock()

	return errors.Join(errs...)
}

func (p *sHTTPAdapter) Consume(pCtx context.Context) (net_message.IMessage, error) {
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
