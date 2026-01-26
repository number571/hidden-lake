package client

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	friend "github.com/number571/hidden-lake/pkg/api/kernel/client/dto"
	"github.com/number571/hidden-lake/pkg/network/request"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) Friend(pAliasName string, pPubKey asymmetric.IPubKey) *friend.SFriend {
	if pPubKey == nil {
		// request: del friend
		return &friend.SFriend{
			FAliasName: pAliasName,
		}
	}
	// request: add friend
	return &friend.SFriend{
		FAliasName: pAliasName,
		FPublicKey: pPubKey.ToString(),
	}
}

func (p *sBuilder) Request(pReq request.IRequest) *request.SRequest {
	return pReq.(*request.SRequest)
}
