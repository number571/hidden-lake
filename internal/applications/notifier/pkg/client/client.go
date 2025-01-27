package client

import (
	"bytes"
	"context"
	"math"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/crypto/random"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fSettings  ISettings
	fBuilder   IBuilder
	fRequester IRequester
}

func NewClient(pSettings ISettings, pBuilder IBuilder, pRequester IRequester) IClient {
	return &sClient{
		fSettings:  pSettings,
		fBuilder:   pBuilder,
		fRequester: pRequester,
	}
}

func (p *sClient) Initialize(pCtx context.Context, pTargets []string, pBody []byte) ([]byte, error) {
	salt := random.NewRandom().GetBytes(CSaltSize)
	hash := hashing.NewHasher(bytes.Join([][]byte{salt, pBody}, []byte{})).ToBytes()
	powPuzzle := puzzle.NewPoWPuzzle(p.fSettings.GetWorkSizeBits())
	proof := powPuzzle.ProofBytes(hash, p.fSettings.GetPowParallel())
	return hash, p.Redirect(pCtx, pTargets, "", proof, salt, pBody)
}

func (p *sClient) Finalyze(
	pCtx context.Context,
	pTargets []string,
	pProof uint64,
	pSalt []byte,
	pBody []byte,
) error {
	return p.fRequester.Broadcast(pCtx, pTargets, p.fBuilder.Finalyze(pProof, pSalt, pBody))
}

func (p *sClient) Redirect(
	pCtx context.Context,
	pTargets []string,
	pIgnore string,
	pProof uint64,
	pSalt []byte,
	pBody []byte,
) error {
	if r := random.NewRandom(); r.GetBool() {
		return p.Finalyze(pCtx, pTargets, pProof, pSalt, pBody)
	}
	randTarget := []string{getRandomTarget(deleteTarget(pTargets, pIgnore))}
	return p.fRequester.Broadcast(pCtx, randTarget, p.fBuilder.Redirect(pProof, pSalt, pBody))
}

func deleteTarget(pFriends []string, pIgnore string) []string {
	if pIgnore == "" {
		return pFriends
	}
	result := make([]string, 0, len(pFriends))
	for _, alias := range pFriends {
		if alias == pIgnore {
			continue
		}
		result = append(result, alias)
	}
	return result
}

// The equally probable randomness of a set of bits.
func getRandomTarget(pFriends []string) string {
	if len(pFriends) == 0 {
		return ""
	}
	random := random.NewRandom()
	lenTargets := uint64(len(pFriends))
	pow2 := uint64(math.Pow(2, math.Ceil(math.Log2(float64(lenTargets)))))
	for {
		u64 := random.GetUint64() % pow2
		if u64 >= lenTargets {
			continue
		}
		return pFriends[u64]
	}
}
