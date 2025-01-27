package database

import "github.com/number571/go-peer/pkg/crypto/asymmetric"

var (
	_ IRelation = &sRelation{}
)

type sRelation struct {
	fIAm asymmetric.IPubKey
}

func NewRelation(pIAm asymmetric.IPubKey) IRelation {
	return &sRelation{
		fIAm: pIAm,
	}
}

func (p *sRelation) IAm() asymmetric.IPubKey {
	return p.fIAm
}
