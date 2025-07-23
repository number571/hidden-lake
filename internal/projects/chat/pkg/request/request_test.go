package request

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/hidden-lake/internal/projects/chat/pkg/settings"
	"github.com/number571/hidden-lake/pkg/request"
)

func TestGetMessageLimitSize(t *testing.T) {
	t.Parallel()

	if limitSize := GetMessageLimitSize(1000); limitSize != 686 {
		t.Fatal("invalid limit size")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	_ = GetMessageLimitSize(8)
}

func TestBuildRequest(t *testing.T) {
	t.Parallel()

	privKey := ed25519.NewKeyFromSeed(nullSeed)

	msg := "hello, world!"
	req := BuildRequest(chanKey, privKey, msg)

	pubKey, _, ok := ValidateRequest(chanKey, req)
	if !ok {
		t.Fatal("validate request")
	}

	if !pubKey.Equal(privKey.Public()) {
		t.Fatal("invalid public key")
	}
}

func TestBuildFailedRequest1(t *testing.T) {
	t.Parallel()

	privKey := ed25519.NewKeyFromSeed(nullSeed)
	msg := []byte("hello, world!")

	if _, _, ok := ValidateRequest(chanKey, buildFailedRequest1(chanKey, privKey, msg)); ok {
		t.Fatal("success validate failed request (1)")
	}

	if _, _, ok := ValidateRequest(chanKey, buildFailedRequest2(chanKey, privKey, msg)); ok {
		t.Fatal("success validate failed request (2)")
	}

	if _, _, ok := ValidateRequest(chanKey, buildFailedRequest3(chanKey, privKey, msg)); ok {
		t.Fatal("success validate failed request (3)")
	}

	if _, _, ok := ValidateRequest(chanKey, buildFailedRequest4(chanKey, privKey, msg)); ok {
		t.Fatal("success validate failed request (4)")
	}

	if _, _, ok := ValidateRequest(chanKey, buildFailedRequest5(chanKey, privKey, msg)); ok {
		t.Fatal("success validate failed request (5)")
	}
}

func buildFailedRequest1(chanKey asymmetric.IPubKey, privKey ed25519.PrivateKey, body []byte) request.IRequest {
	salt := random.NewRandom().GetBytes(16)
	hash := hashing.NewHMACHasher(
		chanKey.ToBytes(),
		hashing.NewHMACHasher(salt, body).ToBytes(),
	).ToBytes()
	return request.NewRequestBuilder().
		WithHost(settings.CProjectFullName + "?").
		WithHead(map[string]string{
			"pubk": hex.EncodeToString(privKey.Public().(ed25519.PublicKey)),
			"salt": hex.EncodeToString(salt),
			"sign": hex.EncodeToString(ed25519.Sign(privKey, hash)),
		}).
		WithBody(body).
		Build()
}

func buildFailedRequest2(chanKey asymmetric.IPubKey, privKey ed25519.PrivateKey, body []byte) request.IRequest {
	salt := random.NewRandom().GetBytes(16)
	hash := hashing.NewHMACHasher(
		chanKey.ToBytes(),
		hashing.NewHMACHasher(salt, body).ToBytes(),
	).ToBytes()
	return request.NewRequestBuilder().
		WithHost(settings.CProjectFullName).
		WithHead(map[string]string{
			"pubk": hex.EncodeToString(privKey.Public().(ed25519.PublicKey)),
			"salt": hex.EncodeToString(salt),
			"sign": hex.EncodeToString(ed25519.Sign(privKey, hash)),
		}).
		WithBody(bytes.Join([][]byte{body, {0x00}}, []byte{})).
		Build()
}

func buildFailedRequest3(chanKey asymmetric.IPubKey, privKey ed25519.PrivateKey, body []byte) request.IRequest {
	salt := random.NewRandom().GetBytes(16)
	hash := hashing.NewHMACHasher(
		chanKey.ToBytes(),
		hashing.NewHMACHasher(salt, body).ToBytes(),
	).ToBytes()
	return request.NewRequestBuilder().
		WithHost(settings.CProjectFullName).
		WithHead(map[string]string{
			"pubk": hex.EncodeToString(privKey.Public().(ed25519.PublicKey)),
			"sign": hex.EncodeToString(ed25519.Sign(privKey, hash)),
		}).
		WithBody(body).
		Build()
}

func buildFailedRequest4(chanKey asymmetric.IPubKey, privKey ed25519.PrivateKey, body []byte) request.IRequest {
	salt := random.NewRandom().GetBytes(16)
	hash := hashing.NewHMACHasher(
		chanKey.ToBytes(),
		hashing.NewHMACHasher(salt, body).ToBytes(),
	).ToBytes()
	return request.NewRequestBuilder().
		WithHost(settings.CProjectFullName).
		WithHead(map[string]string{
			"pubk": hex.EncodeToString(privKey.Public().(ed25519.PublicKey)),
			"salt": hex.EncodeToString(salt) + "_",
			"sign": hex.EncodeToString(ed25519.Sign(privKey, hash)),
		}).
		WithBody(body).
		Build()
}

func buildFailedRequest5(chanKey asymmetric.IPubKey, privKey ed25519.PrivateKey, body []byte) request.IRequest {
	salt := random.NewRandom().GetBytes(16)
	hash := hashing.NewHMACHasher(
		chanKey.ToBytes(),
		hashing.NewHMACHasher(salt, body).ToBytes(),
	).ToBytes()
	salt[0] ^= 1
	return request.NewRequestBuilder().
		WithHost(settings.CProjectFullName).
		WithHead(map[string]string{
			"pubk": hex.EncodeToString(privKey.Public().(ed25519.PublicKey)),
			"salt": hex.EncodeToString(salt),
			"sign": hex.EncodeToString(ed25519.Sign(privKey, hash)),
		}).
		WithBody(body).
		Build()
}
