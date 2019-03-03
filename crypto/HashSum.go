package crypto

import (
	"crypto/sha256"
)

func HashSum(data []byte) []byte {
    var hashed = sha256.Sum256(data)
    return hashed[:]
}
