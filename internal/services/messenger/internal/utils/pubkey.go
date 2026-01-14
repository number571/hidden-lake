package utils

import (
	"context"
	"errors"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
)

func GetFriendPubKeyByAliasName(
	pCtx context.Context,
	client hlk_client.IClient,
	aliasName string,
) (asymmetric.IPubKey, error) {
	friends, err := client.GetFriends(pCtx)
	if err != nil {
		return nil, errors.Join(ErrGetFriends, err)
	}
	friendPubKey, ok := friends[aliasName]
	if !ok {
		return nil, ErrUndefinedPublicKey
	}
	return friendPubKey, nil
}
