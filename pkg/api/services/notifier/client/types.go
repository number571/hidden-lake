package client

import (
	"context"
)

type IClient interface {
	GetIndex(context.Context) error

	GetNotifyLimit(context.Context) (uint64, error)

	RecvNotification(context.Context, string) (string, []byte, error)
	SendNotification(context.Context, string, []byte) error
}

type IRequester interface {
	GetIndex(context.Context) error

	GetNotifyLimit(context.Context) (uint64, error)

	RecvNotification(context.Context, string) (string, []byte, error)
	SendNotification(context.Context, string, []byte) error
}
