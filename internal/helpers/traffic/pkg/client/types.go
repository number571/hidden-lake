package client

import (
	"context"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/helpers/traffic/pkg/config"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPointer(context.Context) (uint64, error)
	GetHash(context.Context, uint64) (string, error)

	GetMessage(context.Context, string) (net_message.IMessage, error)
	PutMessage(context.Context, net_message.IMessage) error
}

type IBuilder interface {
	PutMessage(net_message.IMessage) string
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPointer(context.Context) (uint64, error)
	GetHash(context.Context, uint64) (string, error)

	GetMessage(context.Context, string) (net_message.IMessage, error)
	PutMessage(context.Context, string) error
}
