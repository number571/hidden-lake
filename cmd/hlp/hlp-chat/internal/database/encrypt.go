package database

import (
	"bytes"
	"errors"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
)

func (p *sDatabase) encryptBytes(pBytes []byte) []byte {
	cipher := symmetric.NewCipher(p.fKey[1])
	encMsg := cipher.EncryptBytes(pBytes)
	mac := hashing.NewHMACHasher(p.fKey[0], encMsg).ToBytes()
	return bytes.Join([][]byte{mac, encMsg}, []byte{})
}

func (p *sDatabase) decryptBytes(pEncBytes []byte) ([]byte, error) {
	if len(pEncBytes) < hashing.CHasherSize+symmetric.CCipherBlockSize {
		return nil, errors.New("invalid ciphertext size") // nolint: err113
	}
	mac := hashing.NewHMACHasher(p.fKey[0], pEncBytes[hashing.CHasherSize:]).ToBytes()
	if !bytes.Equal(mac, pEncBytes[:hashing.CHasherSize]) {
		return nil, errors.New("invalid mac") // nolint: err113
	}
	cipher := symmetric.NewCipher(p.fKey[1])
	return cipher.DecryptBytes(pEncBytes[hashing.CHasherSize:]), nil
}
