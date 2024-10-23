package client

import (
	"context"
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/config"
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

func (p *sClient) GetIndex(pCtx context.Context) (string, error) {
	res, err := p.fRequester.GetIndex(pCtx)
	if err != nil {
		return "", fmt.Errorf("get index (client): %w", err)
	}
	return res, nil
}

func (p *sClient) EncryptMessage(pCtx context.Context, pAliasName string, pPayload payload.IPayload64) (net_message.IMessage, error) {
	res, err := p.fRequester.EncryptMessage(pCtx, pAliasName, pPayload)
	if err != nil {
		return nil, fmt.Errorf("encrypt message (client): %w", err)
	}
	return res, nil
}

func (p *sClient) DecryptMessage(pCtx context.Context, pNetMsg net_message.IMessage) (string, payload.IPayload64, error) {
	aliasName, data, err := p.fRequester.DecryptMessage(pCtx, pNetMsg)
	if err != nil {
		return "", nil, fmt.Errorf("decrypt message (client): %w", err)
	}
	return aliasName, data, nil
}

func (p *sClient) GetPubKey(pCtx context.Context) (asymmetric.IPubKey, error) {
	pubKey, err := p.fRequester.GetPubKey(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get public key (client): %w", err)
	}
	return pubKey, nil
}

func (p *sClient) GetFriends(pCtx context.Context) (map[string]asymmetric.IPubKey, error) {
	res, err := p.fRequester.GetFriends(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get friends (client): %w", err)
	}
	return res, nil
}

func (p *sClient) AddFriend(pCtx context.Context, pAliasName string, pPubKey asymmetric.IPubKey) error {
	if err := p.fRequester.AddFriend(pCtx, p.fBuilder.Friend(pAliasName, pPubKey)); err != nil {
		return fmt.Errorf("add friend (client): %w", err)
	}
	return nil
}

func (p *sClient) DelFriend(pCtx context.Context, pAliasName string) error {
	if err := p.fRequester.DelFriend(pCtx, p.fBuilder.Friend(pAliasName, nil)); err != nil {
		return fmt.Errorf("del friend (client): %w", err)
	}
	return nil
}

func (p *sClient) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	res, err := p.fRequester.GetSettings(pCtx)
	if err != nil {
		return nil, fmt.Errorf("get settings (client): %w", err)
	}
	return res, nil
}
