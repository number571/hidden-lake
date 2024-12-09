package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/go-peer/pkg/payload"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/config"
	hle_settings "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
)

const (
	cHandleIndexTemplate          = "http://" + "%s" + hle_settings.CHandleIndexPath
	cHandleMessageEncryptTemplate = "http://" + "%s" + hle_settings.CHandleMessageEncryptPath
	cHandleMessageDecryptTemplate = "http://" + "%s" + hle_settings.CHandleMessageDecryptPath
	cHandleServicePubKeyTemplate  = "http://" + "%s" + hle_settings.CHandleServicePubKeyPath
	cHandleConfigSettingsTemplate = "http://" + "%s" + hle_settings.CHandleConfigSettingsPath
	cHandleConfigFriendsTemplate  = "http://" + "%s" + hle_settings.CHandleConfigFriendsPath
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHost   string
	fClient *http.Client
	fParams net_message.ISettings
}

func NewRequester(pHost string, pClient *http.Client, pParams net_message.ISettings) IRequester {
	return &sRequester{
		fHost:   pHost,
		fClient: pClient,
		fParams: pParams,
	}
}

func (p *sRequester) GetIndex(pCtx context.Context) (string, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return "", errors.Join(ErrBadRequest, err)
	}

	result := string(res)
	if result != hle_settings.CServiceFullName {
		return "", ErrInvalidTitle
	}

	return result, nil
}

func (p *sRequester) GetFriends(pCtx context.Context) (map[string]asymmetric.IPubKey, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	var vFriends []hls_settings.SFriend
	if err := encoding.DeserializeJSON(res, &vFriends); err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}

	result := make(map[string]asymmetric.IPubKey, len(vFriends))
	for _, friend := range vFriends {
		pubKey := asymmetric.LoadPubKey(friend.FPublicKey)
		if pubKey == nil {
			return nil, ErrInvalidPublicKey
		}
		result[friend.FAliasName] = pubKey
	}

	return result, nil
}

func (p *sRequester) AddFriend(pCtx context.Context, pFriend *hls_settings.SFriend) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		pFriend,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) DelFriend(pCtx context.Context, pFriend *hls_settings.SFriend) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		pFriend,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) EncryptMessage(pCtx context.Context, pAliasName string, pPayload payload.IPayload64) (net_message.IMessage, error) {
	resp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleMessageEncryptTemplate, p.fHost),
		hle_settings.SContainer{
			FAliasName: pAliasName,
			FPldHead:   pPayload.GetHead(),
			FHexData:   encoding.HexEncode(pPayload.GetBody()),
		},
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	msg, err := net_message.LoadMessage(p.fParams, string(resp))
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}

	return msg, nil
}

func (p *sRequester) DecryptMessage(pCtx context.Context, pNetMsg net_message.IMessage) (string, payload.IPayload64, error) {
	resp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleMessageDecryptTemplate, p.fHost),
		pNetMsg.ToString(),
	)
	if err != nil {
		return "", nil, errors.Join(ErrBadRequest, err)
	}

	var result hle_settings.SContainer
	if err := encoding.DeserializeJSON(resp, &result); err != nil {
		return "", nil, errors.Join(ErrDecodeResponse, err)
	}

	data := encoding.HexDecode(result.FHexData)
	if data == nil {
		return "", nil, ErrInvalidHexFormat
	}

	return result.FAliasName, payload.NewPayload64(result.FPldHead, data), nil
}

func (p *sRequester) GetPubKey(pCtx context.Context) (asymmetric.IPubKey, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleServicePubKeyTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	pubKey := asymmetric.LoadPubKey(string(res))
	if pubKey == nil {
		return nil, ErrInvalidPublicKey
	}

	return pubKey, nil
}

func (p *sRequester) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigSettingsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	cfgSettings := new(config.SConfigSettings)
	if err := encoding.DeserializeJSON(res, cfgSettings); err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}

	return cfgSettings, nil
}
