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

func (p *sRequester) PushMessage(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) error {
	if err := p.fHLSClient.SendRequest(pCtx, pAliasName, pRequest); err != nil {
		return errors.Join(ErrPushMessage, err)
	}
	return nil
}
