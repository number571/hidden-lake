package broker

import (
	"context"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/storage/cache"
)

var (
	_ IDataBroker = &sDataBroker{}
)

type sDataBroker struct {
	fMutex    *sync.RWMutex
	fChanSize uint64
	fSubLimit uint64
	fSubChans map[string]chan string
	fCache    cache.ICache
}

func NewDataBroker(pChanSize, pSubLimit uint64) IDataBroker {
	return &sDataBroker{
		fMutex:    &sync.RWMutex{},
		fChanSize: pChanSize,
		fSubLimit: pSubLimit,
		fSubChans: make(map[string]chan string, pSubLimit),
		fCache:    cache.NewLRUCache(pChanSize),
	}
}

func (p *sDataBroker) Produce(pData interface{}) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	uniqID := random.NewRandom().GetString(32)
	p.fCache.Set(uniqID, pData)

	for id, ch := range p.fSubChans {
		select {
		case ch <- uniqID:
		default:
			close(ch)
			delete(p.fSubChans, id)
		}
	}
}

func (p *sDataBroker) Register(pID string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if uint64(len(p.fSubChans)) >= p.fSubLimit {
		return ErrLimitSubscribers
	}
	if _, ok := p.fSubChans[pID]; !ok {
		p.fSubChans[pID] = make(chan string, p.fChanSize)
	}

	return nil
}

func (p *sDataBroker) Consume(pCtx context.Context, pID string) (interface{}, error) {
	ch, ok := p.getSubChannel(pID)
	if !ok {
		return nil, ErrNotRegistered
	}
	select {
	case <-pCtx.Done():
		return nil, pCtx.Err()
	case x := <-ch:
		v, ok := p.fCache.Get(x)
		if !ok {
			return nil, ErrValutNotFound
		}
		return v, nil
	}
}

func (p *sDataBroker) CountSubscribers() uint64 {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return uint64(len(p.fSubChans))
}

func (p *sDataBroker) getSubChannel(pID string) (chan string, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	ch, ok := p.fSubChans[pID]
	return ch, ok
}
