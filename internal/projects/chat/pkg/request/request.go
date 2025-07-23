package request

import (
	"crypto/ed25519"
	"encoding/hex"
	"unicode"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/hidden-lake/internal/projects/chat/pkg/settings"
	"github.com/number571/hidden-lake/pkg/request"
)

func BuildRequest(chanKey asymmetric.IPubKey, privKey ed25519.PrivateKey, body string) request.IRequest {
	salt := random.NewRandom().GetBytes(16)
	hash := hashing.NewHMACHasher(
		chanKey.ToBytes(),
		hashing.NewHMACHasher(salt, []byte(body)).ToBytes(),
	).ToBytes()
	return request.NewRequestBuilder().
		WithHost(settings.CProjectFullName).
		WithHead(map[string]string{
			"pubk": hex.EncodeToString(privKey.Public().(ed25519.PublicKey)),
			"salt": hex.EncodeToString(salt),
			"sign": hex.EncodeToString(ed25519.Sign(privKey, hash)),
		}).
		WithBody([]byte(body)).
		Build()
}

func ValidateRequest(chanKey asymmetric.IPubKey, req request.IRequest) (ed25519.PublicKey, []byte, bool) {
	if req.GetHost() != settings.CProjectFullName {
		return nil, nil, false
	}

	if HasNotGraphicCharacters(string(req.GetBody())) {
		return nil, nil, false
	}

	head := req.GetHead()
	pubkHex, ok1 := head["pubk"]
	saltHex, ok2 := head["salt"]
	signHex, ok3 := head["sign"]
	if !ok1 || !ok2 || !ok3 {
		return nil, nil, false
	}

	pubk, err1 := hex.DecodeString(pubkHex)
	salt, err2 := hex.DecodeString(saltHex)
	sign, err3 := hex.DecodeString(signHex)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, nil, false
	}

	pubKey := ed25519.PublicKey(pubk)
	hash := hashing.NewHMACHasher(
		chanKey.ToBytes(),
		hashing.NewHMACHasher(salt, req.GetBody()).ToBytes(),
	).ToBytes()

	return pubKey, hash, ed25519.Verify(pubKey, hash, sign)
}

func HasNotGraphicCharacters(pS string) bool {
	for _, c := range pS {
		if !unicode.IsGraphic(c) {
			return true
		}
	}
	return false
}
