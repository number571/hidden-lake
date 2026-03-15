package broker

import (
	"sync"
)

var (
	_ IDataBroker = &sDataBroker{}
)

type sDataBroker struct {
	fMutex       *sync.RWMutex
	fChanSize    uint64
	fSubLimit    uint64
	fSubscribers map[string]chan interface{}
}

func NewDataBroker(pChanSize, pSubLimit uint64) IDataBroker {
	return &sDataBroker{
		fMutex:       &sync.RWMutex{},
		fChanSize:    pChanSize,
		fSubLimit:    pSubLimit,
		fSubscribers: make(map[string]chan interface{}, pSubLimit),
	}
}

func (p *sDataBroker) Produce(pData interface{}) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	for id, ch := range p.fSubscribers {
		select {
		case ch <- pData:
		default:
			close(ch)
			delete(p.fSubscribers, id)
		}
	}
}

func (p *sDataBroker) Consume(pID string) <-chan interface{} {
	if ch, ok := p.tryGetChannel(pID); ok {
		return ch
	}

	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if uint64(len(p.fSubscribers)) >= p.fSubLimit {
		ch := make(chan interface{}, 1)
		close(ch)
		return ch
	}

	ch, ok := p.fSubscribers[pID]
	if !ok {
		ch = make(chan interface{}, p.fChanSize)
		p.fSubscribers[pID] = ch
	}

	return ch
}

func (p *sDataBroker) CountSubscribers() uint64 {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	return uint64(len(p.fSubscribers))
}

func (p *sDataBroker) tryGetChannel(pID string) (chan interface{}, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	ch, ok := p.fSubscribers[pID]
	return ch, ok
}
