package crypto

import (
    "crypto/aes"
    "crypto/cipher"
)

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

    return PKCS5UnPadding(data), nil
}

func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
