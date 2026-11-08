package client

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/hidden-lake/pkg/api/adapters/http/config"
)

type IClient interface {
	GetIndex(pCtx context.Context, pProtoScheme string) error
	GetSettings(pCtx context.Context) (config.IConfigSettings, error)

	GetOnlines(pCtx context.Context) ([]string, error)
	DelOnline(pCtx context.Context, pAddr string) error

	GetConnections(pCtx context.Context) ([]string, error)
	AddConnection(pCtx context.Context, pAddr string) error
	DelConnection(pCtx context.Context, pAddr string) error

	ProduceMessage(pCtx context.Context, pMessage layer1.IMessage) error
	ConsumeMessage(pCtx context.Context, pSID string) (layer1.IMessage, error)
}
