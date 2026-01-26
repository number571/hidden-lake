package proc

import (
	"sort"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	friend "github.com/number571/hidden-lake/pkg/api/kernel/client/dto"
)

func FriendsMapToList(pFriendsMap map[string]asymmetric.IPubKey) []friend.SFriend {
	listFriends := make([]friend.SFriend, 0, len(pFriendsMap))
	for name, pubKey := range pFriendsMap {
		listFriends = append(listFriends, friend.SFriend{
			FAliasName: name,
			FPublicKey: pubKey.ToString(),
		})
	}
	sort.Slice(listFriends, func(i, j int) bool {
		return listFriends[i].FAliasName < listFriends[j].FAliasName
	})
	return listFriends
}
