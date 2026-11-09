package client

import (
	"context"
	"time"

	message "github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
)

type IClient interface {
	GetIndex(pCtx context.Context) error

	GetMsgBoxSize(pCtx context.Context, pFriend string) (uint64, error)
	ListenMsgbox(pCtx context.Context, pSID string) ([2]string, message.IMessageBody, error)

	GetPushLimit(pCtx context.Context) (uint64, error)
	PushMessageToBox(pCtx context.Context, pFriend string, pUID string, pData []byte) (uint64, time.Time, error)

	LoadMessageFromBoxByUID(pCtx context.Context, pFriend string, pUID string) (uint64, message.IMessage, error)
	LoadMessageFromBoxByEnum(pCtx context.Context, pFriend string, pNID uint64) (message.IMessage, error)
}
