package client

import (
	"context"

	"github.com/number571/hidden-lake/internal/services/messenger/pkg/client/message"
)

type IClient interface {
	GetIndex(context.Context) (string, error)

	GetMessageLimit(context.Context) (uint64, error)
	ListenChat(context.Context, string, string) (message.IMessage, error)
	PushMessage(context.Context, string, string) (string, error)
	LoadMessages(context.Context, string, uint64, uint64, bool) ([]message.IMessage, error)
	CountMessages(context.Context, string) (uint64, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)

	GetMessageLimit(context.Context) (uint64, error)
	ListenChat(context.Context, string, string) (message.IMessage, error)
	PushMessage(context.Context, string, string) (string, error)
	LoadMessages(context.Context, string, uint64, uint64, bool) ([]message.IMessage, error)
	CountMessages(context.Context, string) (uint64, error)
}
