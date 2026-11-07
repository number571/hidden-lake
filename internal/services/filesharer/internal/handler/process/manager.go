package process

import (
	"sort"
	"sync"

	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

var (
	_ IDownloadProcessManager = &sDownloadProcessManager{}
)

type sDownloadProcessManager struct {
	fMutex   *sync.RWMutex
	fCounter uint64
	fMap     map[string]*sDownloadProcessStatic
}

type sDownloadProcessStatic struct {
	fCancel  func()
	fProcess *dto.IDownloadProcess
}

func NewDownloadProcessesManager() IDownloadProcessManager {
	return &sDownloadProcessManager{
		fMutex: &sync.RWMutex{},
		fMap:   make(map[string]*sDownloadProcessStatic, 256),
	}
}

func (p *sDownloadProcessManager) Update(k dto.IDownloadProcessKey, u [2]uint64) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	sk := k.ToString()
	v, ok := p.fMap[sk]
	if !ok {
		return false
	}

	proc := dto.NewDownloadProcess(
		k,
		dto.NewDownloadProcessValue(p.fCounter, u[0], u[1]),
	)
	v.fProcess = &proc

	return true
}

func (p *sDownloadProcessManager) Get(k dto.IDownloadProcessKey) (dto.IDownloadProcess, bool) {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	sk := k.ToString()
	v, ok := p.fMap[sk]
	if !ok {
		return nil, false
	}

	return *v.fProcess, true
}

func (p *sDownloadProcessManager) GetKeys() []dto.IDownloadProcessKey {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	processes := make([]dto.IDownloadProcess, 0, len(p.fMap))
	for _, v := range p.fMap {
		processes = append(processes, *v.fProcess)
	}

	sort.Slice(
		processes,
		func(i, j int) bool { return processes[i].GetIncIndex() < processes[j].GetIncIndex() },
	)

	result := make([]dto.IDownloadProcessKey, 0, len(processes))
	for _, v := range processes {
		result = append(result, v)
	}
	return result
}

func (p *sDownloadProcessManager) TryLock(k dto.IDownloadProcessKey, c func()) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	sk := k.ToString()
	if _, ok := p.fMap[sk]; ok {
		return false
	}

	proc := dto.NewDownloadProcess(
		k,
		dto.NewDownloadProcessValue(p.fCounter, 0, 0),
	)

	p.fCounter++
	p.fMap[sk] = &sDownloadProcessStatic{
		fCancel:  c,
		fProcess: &proc,
	}

	return true
}

func (p *sDownloadProcessManager) Unlock(k dto.IDownloadProcessKey) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	sk := k.ToString()
	val, ok := p.fMap[sk]
	if !ok {
		return false
	}

	val.fCancel()
	delete(p.fMap, sk)

	return true
}
