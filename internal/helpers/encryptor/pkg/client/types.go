package client

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/config"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	EncryptMessage(context.Context, asymmetric.IKEncPubKey, payload.IPayload64) (net_message.IMessage, error)
	DecryptMessage(context.Context, net_message.IMessage) (asymmetric.ISignPubKey, payload.IPayload64, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	EncryptMessage(context.Context, asymmetric.IKEncPubKey, payload.IPayload64) (net_message.IMessage, error)
	DecryptMessage(context.Context, net_message.IMessage) (asymmetric.ISignPubKey, payload.IPayload64, error)
}
