package client

import (
	"context"
	"fmt"

	"github.com/number571/hidden-lake/internal/service/pkg/config"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fRequester IRequester
}

func NewClient(pRequester IRequester) IClient {
	return &sClient{
		fRequester: pRequester,
	}
}

func (p *sClient) GetIndex(pCtx context.Context) (string, error) {
	res, err := p.fRequester.GetIndex(pCtx)
	if err != nil {
		return "", fmt.Errorf("get index (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	res, err := p.fRequester.GetSettings(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get settings (client): %w", err)
	}
	return res, nil
}

func (p *sClient) GetOnlines(pCtx context.Context) ([]string, error) {
	res, err := p.fRequester.GetOnlines(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get onlines (client): %w", err)
	}
	return res, nil
}

func (p *sClient) DelOnline(pCtx context.Context, pConnect string) error {
	if err := p.fRequester.DelOnline(pCtx, pConnect); err != nil {
		return fmt.Errorf("del online (client): %w", err)
	}
	return nil
}

func (p *sClient) GetConnections(pCtx context.Context) ([]string, error) {
	res, err := p.fRequester.GetConnections(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get connections (client): %w", err)
	}
	return res, nil
}

func (p *sClient) AddConnection(pCtx context.Context, pConnect string) error {
	if err := p.fRequester.AddConnection(pCtx, pConnect); err != nil {
		return fmt.Errorf("add connection (client): %w", err)
	}
	return nil
}

func (p *sClient) DelConnection(pCtx context.Context, pConnect string) error {
	if err := p.fRequester.DelConnection(pCtx, pConnect); err != nil {
		return fmt.Errorf("del connection (client): %w", err)
	}
	return nil
}
