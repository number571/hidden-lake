package database

import "github.com/number571/go-peer/pkg/crypto/asymmetric"

var (
	_ IRelation = &sRelation{}
)

type sRelation struct {
	fIAm asymmetric.IPubKey
	fKey string
}

func NewRelation(pIAm asymmetric.IPubKey, pKey string) IRelation {
	return &sRelation{
		fIAm: pIAm,
		fKey: pKey,
	}
}

func (p *sRelation) IAm() asymmetric.IPubKey {
	return p.fIAm
}

func (p *sRelation) Key() string {
	return p.fKey
}
