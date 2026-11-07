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

func (p *sClient) PushMessageToBox(pCtx context.Context, pAliasName string, pUID string, pBody []byte, pRemote bool) (time.Time, error) {
	return p.fRequester.PushMessageToBox(pCtx, pAliasName, pUID, pBody, pRemote)
}

func (p *sClient) LoadMessageFromBoxByUID(pCtx context.Context, pAliasName string, pUID string, pRemote bool) (message.IMessage, error) {
	return p.fRequester.LoadMessageFromBoxByUID(pCtx, pAliasName, pUID, pRemote)
}

func (p *sClient) LoadMessageFromBoxByEnum(pCtx context.Context, pAliasName string, pEnum uint64, pRemote bool) (message.IMessage, error) {
	return p.fRequester.LoadMessageFromBoxByEnum(pCtx, pAliasName, pEnum, pRemote)
}

func (p *sClient) ListenMsgbox(pCtx context.Context, pSubscribeID string) ([2]string, message.IMessageBody, error) {
	return p.fRequester.ListenMsgbox(pCtx, pSubscribeID)
}
