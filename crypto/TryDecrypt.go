package crypto

import (
	"encoding/hex"
)

func TryDecrypt(session_key []byte, data string) int8 {
    decoded, err := hex.DecodeString(data)
    if err != nil { return 1 }
    _, err = DecryptAES(
        decoded,
        session_key,
    )
    if err != nil { return 2 }
    return 0
}
