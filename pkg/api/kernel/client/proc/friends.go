package proc

import (
	"sort"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
)

func FriendsMapToList(pFriendsMap map[string]layer2.IParticipantKey) [][2]string {
	listFriends := make([][2]string, 0, len(pFriendsMap))
	for name, pKey := range pFriendsMap {
		listFriends = append(listFriends, [2]string{name, pKey.ToString()})
	}
	sort.Slice(listFriends, func(i, j int) bool {
		return listFriends[i][0] < listFriends[j][0]
	})
	return listFriends
}
