package crypto

import (
    "io"
    "bytes"
    "crypto/aes"
    "crypto/rand"
    "crypto/cipher"
    "encoding/hex"
)

func Encrypt(session_key []byte, data string) string {
    result, _ := EncryptAES(
        []byte(data),
        session_key,
    )
    return hex.EncodeToString(result)
}

func EncryptAES(data, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        panic(err)
    }

    blockSize := block.BlockSize()
    data = PKCS5Padding(data, blockSize)

    cipherText := make([]byte, blockSize + len(data))

    iv := cipherText[:blockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }

    mode := cipher.NewCBCEncrypter(block, iv)
    mode.CryptBlocks(cipherText[blockSize:], data)

    return cipherText, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext) % blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}
