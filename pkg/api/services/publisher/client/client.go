package client

import (
	"context"
	"time"

	message "github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
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

func (p *sClient) GetPushLimit(pCtx context.Context) (uint64, error) {
	return p.fRequester.GetPushLimit(pCtx)
}

func (p *sClient) GetMsgBoxSize(pCtx context.Context, pFriend string) (uint64, error) {
	return p.fRequester.GetMsgBoxSize(pCtx, pFriend)
}

func (p *sClient) PushMessageToBox(pCtx context.Context, pAliasName string, pUID string, pBody string) (time.Time, error) {
	return p.fRequester.PushMessageToBox(pCtx, pAliasName, pUID, pBody)
}

func (p *sClient) LoadMessageFromBox(pCtx context.Context, pAliasName string, pUID string) (message.IMessage, error) {
	return p.fRequester.LoadMessageFromBox(pCtx, pAliasName, pUID)
}

func (p *sClient) LoadMessageByEnum(pCtx context.Context, pAliasName string, pIndex uint64) (message.IMessage, error) {
	return p.fRequester.LoadMessageByEnum(pCtx, pAliasName, pIndex)
}

func (p *sClient) LoadMessageByUID(pCtx context.Context, pAliasName string, pUID string) (message.IMessage, error) {
	return p.fRequester.LoadMessageByUID(pCtx, pAliasName, pUID)
}

func (p *sClient) ListenMsgbox(pCtx context.Context, pSubscribeID, pAliasName string) ([2]string, message.IMessageBody, error) {
	return p.fRequester.ListenMsgbox(pCtx, pSubscribeID, pAliasName)
}
