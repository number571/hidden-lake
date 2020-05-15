package utils

import (
	"github.com/number571/gopeer"
)

func RandomString(max uint) string {
	list := gopeer.GenerateRandomBytes(max)
	for i := range list {
		list[i] = list[i]%26 + 'A'
	}
	return string(list)
}
