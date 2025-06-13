package msgdata

import "sync"

var (
	_ IMessageBroker = &sMessageBroker{}
)

type sMessageBroker struct {
	fMapping map[string]*sMapQueue
	fMutexes [3]sync.Mutex
}

type sMapQueue struct {
	fMQueue chan sSubscribeMessage
	fClosed chan struct{}
}

type sSubscribeMessage struct {
	SSubscribe
	SMessage
}

func NewMessageBroker() IMessageBroker {
	return &sMessageBroker{
		fMapping: make(map[string]*sMapQueue, 128),
	}
}

func (p *sMessageBroker) Close(pAddress string) bool {
	p.fMutexes[0].Lock()
	defer p.fMutexes[0].Unlock()

	mq, ok := p.loadAndDeleteMapQueue(pAddress)
	if ok {
		p.close(mq)
	}
	return ok
}

func (p *sMessageBroker) Consume(pAddress string) (SMessage, bool) {
	p.fMutexes[0].Lock()
	mq, ok := p.loadMapQueue(pAddress)
	if ok {
		p.close(mq)
	}
	mq = p.createMapQueue(pAddress)
	p.fMutexes[0].Unlock()
	select {
	case <-mq.fClosed:
		return SMessage{}, false
	case msg := <-mq.fMQueue:
		return msg.SMessage, true
	}
}

func (p *sMessageBroker) Produce(pAddress string, pMsg SMessage) {
	p.fMutexes[1].Lock()
	defer p.fMutexes[1].Unlock()

	mq, ok := p.loadMapQueue(pAddress)
	if !ok {
		return
	}

	p.clear(mq) // only one can produce value

	mq.fMQueue <- sSubscribeMessage{
		SSubscribe: SSubscribe{FAddress: pAddress},
		SMessage:   pMsg,
	}
}

func (p *sMessageBroker) close(mq *sMapQueue) {
	p.clear(mq)
	close(mq.fClosed)
}

func (p *sMessageBroker) clear(mq *sMapQueue) {
	if len(mq.fMQueue) > 0 {
		<-mq.fMQueue
	}
}

func (p *sMessageBroker) loadMapQueue(pAddress string) (*sMapQueue, bool) {
	p.fMutexes[2].Lock()
	defer p.fMutexes[2].Unlock()

	mq, ok := p.fMapping[pAddress]
	return mq, ok
}

func (p *sMessageBroker) createMapQueue(pAddress string) *sMapQueue {
	p.fMutexes[2].Lock()
	defer p.fMutexes[2].Unlock()

	mq := &sMapQueue{
		fMQueue: make(chan sSubscribeMessage, 1),
		fClosed: make(chan struct{}),
	}

	p.fMapping[pAddress] = mq
	return mq
}

func (p *sMessageBroker) loadAndDeleteMapQueue(pAddress string) (*sMapQueue, bool) {
	p.fMutexes[2].Lock()
	defer p.fMutexes[2].Unlock()

	mq, ok := p.fMapping[pAddress]
	if !ok {
		return nil, false
	}

	delete(p.fMapping, pAddress)
	return mq, false
}
