package crypto

import (
    "errors"
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

    return unpaddingPKCS5(data)
}

func unpaddingPKCS5(origData []byte) ([]byte, error) {
    length := len(origData)
    unpadding := int(origData[length-1])

    if length < unpadding {
        return nil, errors.New("length < unpadding")
    }

    return origData[:(length - unpadding)], nil
}
