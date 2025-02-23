package client

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/random"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/hidden-lake/internal/utils/rand"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fBuilder   IBuilder
	fRequester IRequester
}

func NewClient(pBuilder IBuilder, pRequester IRequester) IClient {
	return &sClient{
		fBuilder:   pBuilder,
		fRequester: pRequester,
	}
}

func (p *sClient) Finalyze(
	pCtx context.Context,
	pTargets []string,
	pMsg layer1.IMessage,
) error {
	return p.fRequester.Broadcast(pCtx, pTargets, p.fBuilder.Finalyze(pMsg))
}

func (p *sClient) Redirect(
	pCtx context.Context,
	pTargets []string,
	pMsg layer1.IMessage,
) error {
	if r := random.NewRandom(); r.GetBool() {
		return p.Finalyze(pCtx, pTargets, pMsg)
	}
	randTarget := make([]string, 0, 1)
	if x := getRandomTarget(pTargets); x != "" {
		randTarget = append(randTarget, x)
	}
	return p.fRequester.Broadcast(pCtx, randTarget, p.fBuilder.Redirect(pMsg))
}

func getRandomTarget(pFriends []string) string {
	if len(pFriends) == 0 {
		return ""
	}
	lenFriends := uint64(len(pFriends))
	return pFriends[rand.UniformUint64n(lenFriends)]
}
