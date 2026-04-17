package limiter

import (
	"sync"
	"time"
)

var (
	_ ILimiter = &sLimiter{}
)

type sLimiter struct {
	fUpdatedAt time.Time
	fMtx       *sync.Mutex
	fTokens    float64
	fRate      float64
	fCapacity  float64
}

func NewLimiter(pRate, pCapacity float64) ILimiter {
	return &sLimiter{
		fMtx:      &sync.Mutex{},
		fRate:     pRate,
		fCapacity: pCapacity,
	}
}

func (p *sLimiter) Allow() bool {
	p.fMtx.Lock()
	defer p.fMtx.Unlock()

	duration := time.Since(p.fUpdatedAt).Seconds()
	p.fUpdatedAt = time.Now()

	p.fTokens += (duration * p.fRate)

	if p.fTokens > p.fCapacity {
		p.fTokens = p.fCapacity
	}

	ok := p.fTokens >= 1
	if ok {
		p.fTokens--
	}

	return ok
}
