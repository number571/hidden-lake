package crypto

import (
	"encoding/hex"
)

func Encrypt(session_key []byte, data string) string {
    result, _ := EncryptAES(
        []byte(data),
        session_key,
    )
    return hex.EncodeToString(result)
}
