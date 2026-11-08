package client

import (
	"context"
)

type IClient interface {
	GetIndex(pCtx context.Context) error

	GetNotifyLimit(pCtx context.Context) (uint64, error)

	RecvNotification(pCtx context.Context, pSID string) (string, []byte, error)
	SendNotification(pCtx context.Context, pFriend string, pData []byte) error
}
