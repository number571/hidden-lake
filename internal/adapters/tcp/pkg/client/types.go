package client

import (
	"context"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/service/pkg/config"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetOnlines(context.Context) ([]string, error)
	DelOnline(context.Context, string) error

	GetConnections(context.Context) ([]string, error)
	AddConnection(context.Context, string) error
	DelConnection(context.Context, string) error

	ProduceMessage(context.Context, net_message.IMessage) error
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetOnlines(context.Context) ([]string, error)
	DelOnline(context.Context, string) error

	GetConnections(context.Context) ([]string, error)
	AddConnection(context.Context, string) error
	DelConnection(context.Context, string) error

	ProduceMessage(context.Context, net_message.IMessage) error
}
