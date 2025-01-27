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

func (p *sClient) Notify(
	pCtx context.Context,
	pTargets []string,
	pIgnore string,
	pProof uint64,
	pSalt []byte,
	pBody []byte,
) ([]byte, error) {
	var hash []byte
	if len(pTargets) == 0 {
		return nil, ErrTargetsIsNull
	}
	if pIgnore == "" {
		pSalt = random.NewRandom().GetBytes(CSaltSize)
		hash = hashing.NewHasher(bytes.Join([][]byte{pSalt, pBody}, []byte{})).ToBytes()
		powPuzzle := puzzle.NewPoWPuzzle(p.fSettings.GetDiffBits())
		pProof = powPuzzle.ProofBytes(hash, p.fSettings.GetParallel())
	}
	// send to all friends
	if r := random.NewRandom(); r.GetBool() {
		return hash, p.Finalyze(pCtx, pTargets, pProof, pSalt, pBody)
	}
	// send to one friend
	return hash, p.Redirect(pCtx, pTargets, pIgnore, pProof, pSalt, pBody)
}

func (p *sClient) Finalyze(
	pCtx context.Context,
	pTargets []string,
	pProof uint64,
	pSalt []byte,
	pBody []byte,
) error {
	return p.fRequester.Finalyze(pCtx, pTargets, p.fBuilder.Finalyze(pProof, pSalt, pBody))
}

func (p *sClient) Redirect(
	pCtx context.Context,
	pTargets []string,
	pIgnore string,
	pProof uint64,
	pSalt []byte,
	pBody []byte,
) error {
	randTarget := getRandomTarget(deleteTarget(pTargets, pIgnore))
	return p.fRequester.Redirect(pCtx, randTarget, p.fBuilder.Redirect(pProof, pSalt, pBody))
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

func getRandomTarget(pFriends []string) string {
	// the equally probable randomness of a set of bits
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
