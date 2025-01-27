package client

import (
	"context"
	"errors"

	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHLSClient hls_client.IClient
}

func NewRequester(pHLSClient hls_client.IClient) IRequester {
	return &sRequester{
		fHLSClient: pHLSClient,
	}
}

func (p *sRequester) Finalyze(pCtx context.Context, pAliasNames []string, pRequest hls_request.IRequest) error {
	for _, alias := range pAliasNames {
		if err := p.fHLSClient.SendRequest(pCtx, alias, pRequest); err != nil {
			return errors.Join(ErrBadRequest, err)
		}
	}
	return nil
}

func (p *sRequester) Redirect(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) error {
	if err := p.fHLSClient.SendRequest(pCtx, pAliasName, pRequest); err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}
