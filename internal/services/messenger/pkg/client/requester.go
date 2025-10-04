package client

import (
	"context"
	"errors"

	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hlk_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHlkClient hlk_client.IClient
}

func NewRequester(pHlkClient hlk_client.IClient) IRequester {
	return &sRequester{
		fHlkClient: pHlkClient,
	}
}

func (p *sRequester) PushMessage(pCtx context.Context, pAliasName string, pRequest hlk_request.IRequest) error {
	if err := p.fHlkClient.SendRequest(pCtx, pAliasName, pRequest); err != nil {
		return errors.Join(ErrPushMessage, err)
	}
	return nil
}
