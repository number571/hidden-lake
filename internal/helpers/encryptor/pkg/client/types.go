package client

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/config"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	GetFriends(context.Context) (map[string]asymmetric.IPubKey, error)
	AddFriend(context.Context, string, asymmetric.IPubKey) error
	DelFriend(context.Context, string) error

	EncryptMessage(context.Context, string, payload.IPayload64) (net_message.IMessage, error)
	DecryptMessage(context.Context, net_message.IMessage) (string, payload.IPayload64, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	GetFriends(context.Context) (map[string]asymmetric.IPubKey, error)
	AddFriend(context.Context, *hls_settings.SFriend) error
	DelFriend(context.Context, *hls_settings.SFriend) error

	EncryptMessage(context.Context, string, payload.IPayload64) (net_message.IMessage, error)
	DecryptMessage(context.Context, net_message.IMessage) (string, payload.IPayload64, error)
}

type IBuilder interface {
	Friend(string, asymmetric.IPubKey) *hls_settings.SFriend
}
