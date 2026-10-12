package client

import (
	"context"
	"time"

	message "github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
)

type IClient interface {
	GetIndex(context.Context) error

	GetNotifyLimit(context.Context) (uint64, error)
	GetMessageLimit(context.Context) (uint64, error)

	GetChatSize(context.Context, string) (uint64, error)

	ListenChan(context.Context, string, string) (message.IMessageBody, error)
	LoadMessage(context.Context, string, uint64) (message.IMessage, error)

	PushMessage(context.Context, string, string, string) (time.Time, error)
	SendNotification(context.Context, string, string) (time.Time, error)
}

type IRequester interface {
	GetIndex(context.Context) error

	GetNotifyLimit(context.Context) (uint64, error)
	GetMessageLimit(context.Context) (uint64, error)

	GetChatSize(context.Context, string) (uint64, error)

	ListenChan(context.Context, string, string) (message.IMessageBody, error)
	LoadMessage(context.Context, string, uint64) (message.IMessage, error)

	PushMessage(context.Context, string, string, string) (time.Time, error)
	SendNotification(context.Context, string, string) (time.Time, error)
}
