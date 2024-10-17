package database

import "github.com/number571/go-peer/pkg/crypto/asymmetric"

var (
	_ IRelation = &sRelation{}
)

type sRelation struct {
	fIAm    asymmetric.IPubKeyChain
	fFriend asymmetric.IPubKeyChain
}

func NewRelation(pIAm, pFriend asymmetric.IPubKeyChain) IRelation {
	return &sRelation{
		fIAm:    pIAm,
		fFriend: pFriend,
	}
}

func (p *sRelation) IAm() asymmetric.IPubKeyChain {
	return p.fIAm
}

func (p *sRelation) Friend() asymmetric.IPubKeyChain {
	return p.fFriend
}
