package client

import (
	"context"
	"time"

	message "github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
)

type IClient interface {
	GetIndex(context.Context) error

	GetPushLimit(context.Context) (uint64, error)
	GetMsgBoxSize(context.Context, string) (uint64, error)

	ListenMsgbox(context.Context, string, string) ([2]string, message.IMessageBody, error)

	LoadMessageByUID(context.Context, string, string) (message.IMessage, error)
	LoadMessageByEnum(context.Context, string, uint64) (message.IMessage, error)

	PushMessageToBox(context.Context, string, string, string) (time.Time, error)
	LoadMessageFromBox(context.Context, string, string) (message.IMessage, error)
}

type IRequester interface {
	GetIndex(context.Context) error

	GetPushLimit(context.Context) (uint64, error)
	GetMsgBoxSize(context.Context, string) (uint64, error)

	ListenMsgbox(context.Context, string, string) ([2]string, message.IMessageBody, error)

	LoadMessageByUID(context.Context, string, string) (message.IMessage, error)
	LoadMessageByEnum(context.Context, string, uint64) (message.IMessage, error)

	PushMessageToBox(context.Context, string, string, string) (time.Time, error)
	LoadMessageFromBox(context.Context, string, string) (message.IMessage, error)
}
