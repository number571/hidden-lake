//go:build !symmetric

package utils

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
)

func LoadParticipantKey(pFriendKey string) layer2.IParticipantKey {
	return asymmetric.LoadPubKey(pFriendKey)
}
