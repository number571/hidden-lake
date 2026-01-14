package message

import "sync"

var (
	_ IMessageBroker = &sMessageBroker{}
)

type sMessageBroker struct {
	fMutex       *sync.RWMutex
	fSubscribers map[string]chan IMessageContainer
}

func NewMessageBroker() IMessageBroker {
	return &sMessageBroker{
		fMutex:       &sync.RWMutex{},
		fSubscribers: make(map[string]chan IMessageContainer, 256),
	}
}

func (p *sMessageBroker) Produce(pFriend string, pMessage IMessage) {
	c := newMessageContainer(pFriend, pMessage)
	for id, ch := range p.fSubscribers {
		select {
		case ch <- c:
		default:
			p.fMutex.Lock()
			delete(p.fSubscribers, id)
			p.fMutex.Unlock()
		}
	}
}

func (p *sMessageBroker) Consume(pID string) chan IMessageContainer {
	if ch, ok := p.tryGetChannel(pID); ok {
		return ch
	}

	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	ch, ok := p.fSubscribers[pID]
	if !ok {
		ch = make(chan IMessageContainer, 32)
		p.fSubscribers[pID] = ch
	}

	return ch
}

func (p *sMessageBroker) tryGetChannel(pID string) (chan IMessageContainer, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	ch, ok := p.fSubscribers[pID]
	return ch, ok
}
