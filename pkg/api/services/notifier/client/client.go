package client

import (
	"context"
	"time"

	message "github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
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

func (p *sClient) GetNotifyLimit(pCtx context.Context) (uint64, error) {
	return p.fRequester.GetNotifyLimit(pCtx)
}

func (p *sClient) GetMessageLimit(pCtx context.Context) (uint64, error) {
	return p.fRequester.GetMessageLimit(pCtx)
}

func (p *sClient) SendNotification(pCtx context.Context, pAliasName string, pBody string) (time.Time, error) {
	return p.fRequester.SendNotification(pCtx, pAliasName, pBody)
}

func (p *sClient) SendMessageToBox(pCtx context.Context, pAliasName string, pUID string, pBody string) (time.Time, error) {
	return p.fRequester.SendMessageToBox(pCtx, pAliasName, pUID, pBody)
}

func (p *sClient) LoadMessageByEnum(pCtx context.Context, pAliasName string, pIndex uint64) (message.IMessage, error) {
	return p.fRequester.LoadMessageByEnum(pCtx, pAliasName, pIndex)
}

func (p *sClient) LoadMessageByUID(pCtx context.Context, pAliasName string, pUID string) (message.IMessage, error) {
	return p.fRequester.LoadMessageByUID(pCtx, pAliasName, pUID)
}

func (p *sClient) GetMsgBoxSize(pCtx context.Context, pAliasName string) (uint64, error) {
	return p.fRequester.GetMsgBoxSize(pCtx, pAliasName)
}

func (p *sClient) ListenChan(pCtx context.Context, pSubscribeID, pAliasName string) (string, message.IMessageBody, error) {
	return p.fRequester.ListenChan(pCtx, pSubscribeID, pAliasName)
}
