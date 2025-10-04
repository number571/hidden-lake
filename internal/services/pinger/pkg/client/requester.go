package client

import (
	"context"
	"errors"
	"net/http"

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

func (p *sRequester) Ping(pCtx context.Context, pAliasName string, pRequest hlk_request.IRequest) error {
	resp, err := p.fHlkClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return ErrDecodeResponse
	}

	return nil
}
