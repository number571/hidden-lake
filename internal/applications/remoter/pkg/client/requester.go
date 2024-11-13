package client

import (
	"context"
	"errors"
	"net/http"

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

func (p *sRequester) Exec(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) ([]byte, error) {
	resp, err := p.fHLSClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return nil, ErrDecodeResponse
	}

	return resp.GetBody(), nil
}
