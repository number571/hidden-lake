package client

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"

	friend "github.com/number571/hidden-lake/pkg/api/kernel/client/dto"
	"github.com/number571/hidden-lake/pkg/api/kernel/config"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

type IClient interface {
	GetIndex(context.Context) error
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
	GetIndex(context.Context) error
	GetSettings(context.Context) (config.IConfigSettings, error)

	GetPubKey(context.Context) (asymmetric.IPubKey, error)

	GetOnlines(context.Context) ([]string, error)
	DelOnline(context.Context, string) error

	GetFriends(context.Context) (map[string]asymmetric.IPubKey, error)
	AddFriend(context.Context, *friend.SFriend) error
	DelFriend(context.Context, *friend.SFriend) error

	GetConnections(context.Context) ([]string, error)
	AddConnection(context.Context, string) error
	DelConnection(context.Context, string) error

	SendRequest(context.Context, string, *request.SRequest) error
	FetchRequest(context.Context, string, *request.SRequest) (response.IResponse, error)
}

type IBuilder interface {
	Request(request.IRequest) *request.SRequest
	Friend(string, asymmetric.IPubKey) *friend.SFriend
}
