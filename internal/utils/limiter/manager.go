package limiter

import "sync"

var (
	_ ILimitManager = &sLimitManager{}
)

type sLimitManager struct {
	fMtx      *sync.RWMutex
	fLimiters map[string]ILimiter
	fRate     float64
	fCapacity float64
}

func NewLimitManager(pRate, pCapacity float64) ILimitManager {
	return &sLimitManager{
		fMtx:      &sync.RWMutex{},
		fLimiters: make(map[string]ILimiter, 256),
		fRate:     pRate,
		fCapacity: pCapacity,
	}
}

func (p *sLimitManager) Get(pKey string) ILimiter {
	p.fMtx.Lock()
	defer p.fMtx.Unlock()

	v, ok := p.fLimiters[pKey]
	if !ok {
		v = NewLimiter(p.fRate, p.fCapacity)
		p.fLimiters[pKey] = v
	}

	return v
}
