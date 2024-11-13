package client

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"

	"github.com/number571/hidden-lake/internal/service/pkg/config"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

type IClient interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	GetOnlines(context.Context) ([]string, error)
	DelOnline(context.Context, string) error

	GetFriends(context.Context) (map[string]asymmetric.IPubKey, error)
	AddFriend(context.Context, string, asymmetric.IPubKey) error
	DelFriend(context.Context, string) error

	GetConnections(context.Context) ([]string, error)
	AddConnection(context.Context, string) error
	DelConnection(context.Context, string) error

	SendRequest(context.Context, string, request.IRequest) error
	FetchRequest(context.Context, string, request.IRequest) (response.IResponse, error)
}

type IRequester interface {
	GetIndex(context.Context) (string, error)
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	GetOnlines(context.Context) ([]string, error)
	DelOnline(context.Context, string) error

	GetFriends(context.Context) (map[string]asymmetric.IPubKey, error)
	AddFriend(context.Context, *pkg_settings.SFriend) error
	DelFriend(context.Context, *pkg_settings.SFriend) error

	GetConnections(context.Context) ([]string, error)
	AddConnection(context.Context, string) error
	DelConnection(context.Context, string) error

	SendRequest(context.Context, *pkg_settings.SRequest) error
	FetchRequest(context.Context, *pkg_settings.SRequest) (response.IResponse, error)
}

type IBuilder interface {
	Request(string, request.IRequest) *pkg_settings.SRequest
	Friend(string, asymmetric.IPubKey) *pkg_settings.SFriend
}
