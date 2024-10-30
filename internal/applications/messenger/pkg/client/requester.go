package client

import (
	"context"
	"errors"

	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_request "github.com/number571/hidden-lake/internal/service/pkg/request"
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

func (p *sRequester) PingMessage(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) error {
	if _, err := p.fHLSClient.FetchRequest(pCtx, pAliasName, pRequest); err != nil {
		return errors.Join(ErrPingMessage, err)
	}
	return nil
}

func (p *sRequester) PushMessage(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) error {
	if err := p.fHLSClient.BroadcastRequest(pCtx, pAliasName, pRequest); err != nil {
		return errors.Join(ErrPushMessage, err)
	}
	return nil
}
