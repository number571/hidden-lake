package client

import (
	"context"
	"time"

	message "github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
)

type IClient interface {
	GetIndex(context.Context) error

	GetMsgBoxSize(context.Context, string) (uint64, error)
	ListenMsgbox(context.Context, string) ([2]string, message.IMessageBody, error)

	GetPushLimit(context.Context) (uint64, error)
	PushMessageToBox(context.Context, string, string, []byte, bool) (time.Time, error)

	LoadMessageFromBoxByUID(context.Context, string, string, bool) (message.IMessage, error)
	LoadMessageFromBoxByEnum(context.Context, string, uint64, bool) (message.IMessage, error)
}

type IRequester interface {
	GetIndex(context.Context) error

	GetMsgBoxSize(context.Context, string) (uint64, error)
	ListenMsgbox(context.Context, string) ([2]string, message.IMessageBody, error)

	GetPushLimit(context.Context) (uint64, error)
	PushMessageToBox(context.Context, string, string, []byte, bool) (time.Time, error)

	LoadMessageFromBoxByUID(context.Context, string, string, bool) (message.IMessage, error)
	LoadMessageFromBoxByEnum(context.Context, string, uint64, bool) (message.IMessage, error)
}
