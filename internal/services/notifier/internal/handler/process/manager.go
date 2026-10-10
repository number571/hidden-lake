package process

import (
	"sync"
)

var (
	_ INotifyProcessManager = &sNotifyProcessManager{}
)

type sNotifyProcessManager struct {
	fMutex *sync.RWMutex
	fMap   map[string]struct{}
}

func NewNotifyProcessManager() INotifyProcessManager {
	return &sNotifyProcessManager{
		fMutex: &sync.RWMutex{},
		fMap:   make(map[string]struct{}, 256),
	}
}

func (p *sNotifyProcessManager) TryLock(k string) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if _, ok := p.fMap[k]; ok {
		return false
	}

	p.fMap[k] = struct{}{}
	return true
}

func (p *sNotifyProcessManager) Unlock(k string) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if _, ok := p.fMap[k]; !ok {
		return false
	}

	delete(p.fMap, k)
	return true
}
