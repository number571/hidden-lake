package database

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	keyCountMsgsTemplate = "[%X].messages[%X].count"
	keyGetMsgTemplate    = "[%X].messages[%X].get(%d)"
)

func (p *sDatabase) keyCountMsgs(pChannel asymmetric.IPubKey) []byte {
	return []byte(fmt.Sprintf(
		keyCountMsgsTemplate,
		p.fKey[2],
		pChannel.GetHasher().ToBytes(),
	))
}

func (p *sDatabase) keyGetMsg(pChannel asymmetric.IPubKey, i uint64) []byte {
	return []byte(fmt.Sprintf(
		keyGetMsgTemplate,
		p.fKey[2],
		pChannel.GetHasher().ToBytes(),
		i,
	))
}
