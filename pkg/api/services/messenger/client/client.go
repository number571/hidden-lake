package client

import (
	"context"

	message "github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
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
	return p.fRequester.GetIndex(pCtx)
}

func (p *sClient) GetMessageLimit(pCtx context.Context) (uint64, error) {
	return p.fRequester.GetMessageLimit(pCtx)
}

func (p *sClient) PushMessage(pCtx context.Context, pAliasName string, pBody string) (string, error) {
	return p.fRequester.PushMessage(pCtx, pAliasName, pBody)
}

func (p *sClient) LoadMessages(pCtx context.Context, pAliasName string, pStart uint64, pCount uint64, pDesc bool) ([]message.IMessage, error) {
	return p.fRequester.LoadMessages(pCtx, pAliasName, pStart, pCount, pDesc)
}

func (p *sClient) CountMessages(pCtx context.Context, pAliasName string) (uint64, error) {
	return p.fRequester.CountMessages(pCtx, pAliasName)
}

func (p *sClient) ListenChat(pCtx context.Context, pSubscribeID, pAliasName string) (message.IMessage, error) {
	return p.fRequester.ListenChat(pCtx, pSubscribeID, pAliasName)
}
