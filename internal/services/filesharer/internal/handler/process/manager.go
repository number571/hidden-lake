package process

import (
	"fmt"
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

func NewDownloadProcessesMap() IDownloadProcessManager {
	return &sDownloadProcessManager{
		fMutex: &sync.RWMutex{},
		fMap:   make(map[string]*sDownloadProcessStatic, 256),
	}
}

func (p *sDownloadProcessManager) Update(k dto.IDownloadProcessKey, u [2]uint64) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	sk := getDownloadProcessKey(k)
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

	sk := getDownloadProcessKey(k)
	v, ok := p.fMap[sk]
	if !ok {
		return nil, false
	}

	return *v.fProcess, true
}

func (p *sDownloadProcessManager) GetList() []dto.IDownloadProcess {
	p.fMutex.RLock()
	defer p.fMutex.RUnlock()

	result := make([]dto.IDownloadProcess, 0, len(p.fMap))
	for _, v := range p.fMap {
		result = append(result, *v.fProcess)
	}

	sort.Slice(result, func(i, j int) bool { return result[i].GetIncIndex() < result[j].GetIncIndex() })

	return result
}

func (p *sDownloadProcessManager) TryLock(k dto.IDownloadProcessKey, c func()) bool {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	sk := getDownloadProcessKey(k)
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

	sk := getDownloadProcessKey(k)
	val, ok := p.fMap[sk]
	if !ok {
		return false
	}

	val.fCancel()
	delete(p.fMap, sk)

	return true
}

func getDownloadProcessKey(k dto.IDownloadProcessKey) string {
	return fmt.Sprintf("friend=%s&name=%s&personal=%t", k.GetFriend(), k.GetFileName(), k.GetIsPersonal())
}
