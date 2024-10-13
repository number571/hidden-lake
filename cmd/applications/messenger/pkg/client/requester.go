package client

import (
	"context"

	"github.com/number571/go-peer/pkg/utils"
	hls_client "github.com/number571/hidden-lake/cmd/service/pkg/client"
	hls_request "github.com/number571/hidden-lake/cmd/service/pkg/request"
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
	if err := p.fHLSClient.BroadcastRequest(pCtx, pAliasName, pRequest); err != nil {
		return utils.MergeErrors(ErrPushMessage, err)
	}
	return nil
}
