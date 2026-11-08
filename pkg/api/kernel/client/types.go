package client

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"

	"github.com/number571/hidden-lake/pkg/api/kernel/client/scheme"
	"github.com/number571/hidden-lake/pkg/api/kernel/config"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

type IClient interface {
	GetIndex(pCtx context.Context) error
	GetSettings(pCtx context.Context) (config.IConfigSettings, error)

	GetOnlines(pCtx context.Context) ([]string, error)
	DelOnline(pCtx context.Context, pAddr string) error

	GetFriends(pCtx context.Context, pSchemeType scheme.ISchemeType) (map[string]layer2.IParticipantKey, error)
	AddFriend(pCtx context.Context, pFriend string, pKey layer2.IParticipantKey) error
	DelFriend(pCtx context.Context, pFriend string) error

	GetConnections(pCtx context.Context) ([]string, error)
	AddConnection(pCtx context.Context, pAddr string) error
	DelConnection(pCtx context.Context, pAddr string) error

	SendRequest(pCtx context.Context, pFriend string, pRequest request.IRequest) error
	FetchRequest(pCtx context.Context, pFriend string, pRequest request.IRequest) (response.IResponse, error)
}
