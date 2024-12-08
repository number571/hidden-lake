package adapters

import (
	"context"
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/anonymity/adapters"
	net_message "github.com/number571/go-peer/pkg/network/message"
)

var (
	_ IRunnerAdapter = &sMergedRunnerAdapter{}
)

type sMergedRunnerAdapter struct {
	fRunAdapters   []IRunnerAdapter
	fMergedAdapter adapters.IAdapter
}

func NewMergedRunnerAdapter(pRunAdapters ...IRunnerAdapter) IRunnerAdapter {
	iAdapters := make([]adapters.IAdapter, 0, len(pRunAdapters))
	for _, a := range pRunAdapters {
		iAdapters = append(iAdapters, a)
	}
	return &sMergedRunnerAdapter{
		fRunAdapters:   pRunAdapters,
		fMergedAdapter: adapters.NewMergedAdapter(iAdapters...),
	}
}

func (p *sMergedRunnerAdapter) Run(pCtx context.Context) error {
	chCtx, cancel := context.WithCancel(pCtx)
	defer cancel()

	N := len(p.fRunAdapters)

	errs := make([]error, N)
	wg := &sync.WaitGroup{}
	wg.Add(N)

	for i, a := range p.fRunAdapters {
		go func(i int, a IRunnerAdapter) {
			defer func() { wg.Done(); cancel() }()
			errs[i] = a.Run(chCtx)
		}(i, a)
	}

	wg.Wait()

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	default:
		return errors.Join(errs...)
	}
}

func (p *sMergedRunnerAdapter) Produce(pCtx context.Context, pMsg net_message.IMessage) error {
	return p.fMergedAdapter.Produce(pCtx, pMsg)
}

func (p *sMergedRunnerAdapter) Consume(pCtx context.Context) (net_message.IMessage, error) {
	return p.fMergedAdapter.Consume(pCtx)
}
