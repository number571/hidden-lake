package client

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) Friend(pAliasName string, pPubKey asymmetric.IPubKey) *pkg_settings.SFriend {
	if pPubKey == nil {
		// del friend
		return &pkg_settings.SFriend{
			FAliasName: pAliasName,
		}
	}
	// add friend
	return &pkg_settings.SFriend{
		FAliasName: pAliasName,
		FPublicKey: pPubKey.ToString(),
	}
}

func (p *sBuilder) Request(pReceiver string, pReq request.IRequest) *pkg_settings.SRequest {
	return &pkg_settings.SRequest{
		FReceiver: pReceiver,
		FReqData:  pReq.(*request.SRequest),
	}
}
