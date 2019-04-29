package crypto

import (
    "crypto/aes"
    "crypto/cipher"
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

func DecryptAES(data, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    blockSize := block.BlockSize()

    if len(data) < blockSize {
        panic("ciphertext too short")
    }
    
    iv := data[:blockSize]
    data = data[blockSize:]

    if len(data)%blockSize != 0 {
        panic("ciphertext is not a multiple of the block size")
    }

    mode := cipher.NewCBCDecrypter(block, iv)
    mode.CryptBlocks(data, data)

    return PKCS5Unpadding(data), nil
}

func PKCS5Unpadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
