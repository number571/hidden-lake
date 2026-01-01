package adapters

import (
	"context"

	"github.com/number571/go-peer/pkg/anonymity/qb/adapters"
)

var (
	_ IRunnerAdapter = &sRunnerAdapter{}
)

type (
	iRunnerF func(context.Context) error
)

type sRunnerAdapter struct {
	adapters.IAdapter
	fRun iRunnerF
}

func NewRunnerAdapter(pAdapter adapters.IAdapter, pRun iRunnerF) IRunnerAdapter {
	return &sRunnerAdapter{
		IAdapter: pAdapter,
		fRun:     pRun,
	}
}

func (p *sRunnerAdapter) Run(pCtx context.Context) error {
	return p.fRun(pCtx)
}
