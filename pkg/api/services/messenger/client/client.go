package client

import (
	"context"
	"time"

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

func (p *sClient) GetIndex(pCtx context.Context) error {
	return p.fRequester.GetIndex(pCtx)
}

func (p *sClient) GetMessageLimit(pCtx context.Context) (uint64, error) {
	return p.fRequester.GetMessageLimit(pCtx)
}

func (p *sClient) PushMessage(pCtx context.Context, pAliasName string, pBody string) (time.Time, error) {
	return p.fRequester.PushMessage(pCtx, pAliasName, pBody)
}

func (p *sClient) LoadMessage(pCtx context.Context, pAliasName string, pIndex uint64) (message.IMessage, error) {
	return p.fRequester.LoadMessage(pCtx, pAliasName, pIndex)
}

func (p *sClient) GetChatSize(pCtx context.Context, pAliasName string) (uint64, error) {
	return p.fRequester.GetChatSize(pCtx, pAliasName)
}

func (p *sClient) ListenChat(pCtx context.Context, pSubscribeID, pAliasName string) (message.IMessage, error) {
	return p.fRequester.ListenChat(pCtx, pSubscribeID, pAliasName)
}
