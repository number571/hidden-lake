package crypto

import (
	"encoding/hex"
)

func Decrypt(session_key []byte, data string) string {
    decoded, _ := hex.DecodeString(data)
    result, _ := DecryptAES(
        decoded,
        session_key,
    )
    return string(result)
}
