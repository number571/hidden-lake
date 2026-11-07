package process

import (
	"sync"
)

var (
	_ IPublishProcessManager = &sPublishProcessManager{}
)

type sPublishProcessManager struct {
	fMutex *sync.RWMutex
	fMap   map[string]struct{}
}

func NewPublishProcessManager() IPublishProcessManager {
	return &sPublishProcessManager{
		fMutex: &sync.RWMutex{},
		fMap:   make(map[string]struct{}, 256),
	}
}

func (p *sPublishProcessManager) TryLock(k string) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if _, ok := p.fMap[k]; ok {
		return false
	}

	p.fMap[k] = struct{}{}
	return true
}

func (p *sPublishProcessManager) Unlock(k string) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if _, ok := p.fMap[k]; !ok {
		return false
	}

	delete(p.fMap, k)
	return true
}
