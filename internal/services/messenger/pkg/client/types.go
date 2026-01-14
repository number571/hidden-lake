package client

import (
	"context"

	"github.com/number571/hidden-lake/internal/services/messenger/pkg/message"
)

type IClient interface {
	PushMessage(context.Context, string, string) (string, error)
	LoadMessages(context.Context, string, uint64, uint64) ([]message.IMessage, error)
	ListenMessage(context.Context, string, string) (message.IMessage, error)
}

type IRequester interface {
	PushMessage(context.Context, string, string) (string, error)
	LoadMessages(context.Context, string, uint64, uint64) ([]message.IMessage, error)
	ListenMessage(context.Context, string, string) (message.IMessage, error)
}
