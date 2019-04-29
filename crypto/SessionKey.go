package crypto

import (
	"crypto/rand"
)

func SessionKey(max int) []byte {
    var slice []byte = make([]byte, max)
    _, err := rand.Read(slice)
    if err != nil { return nil }
    for max = max - 1; max >= 0; max-- {
        slice[max] = slice[max] % 94 + 33
    }
    return slice
}
