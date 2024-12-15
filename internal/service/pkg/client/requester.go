package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/service/pkg/config"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/pkg/response"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate          = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleConfigSettingsTemplate = "http://" + "%s" + hls_settings.CHandleConfigSettingsPath
	cHandleConfigConnectsTemplate = "http://" + "%s" + hls_settings.CHandleConfigConnectsPath
	cHandleConfigFriendsTemplate  = "http://" + "%s" + hls_settings.CHandleConfigFriendsPath
	cHandleNetworkOnlineTemplate  = "http://" + "%s" + hls_settings.CHandleNetworkOnlinePath
	cHandleNetworkRequestTemplate = "http://" + "%s" + hls_settings.CHandleNetworkRequestPath
	cHandleServicePubKeyTemplate  = "http://" + "%s" + hls_settings.CHandleServicePubKeyPath
)

type sRequester struct {
	fHost   string
	fClient *http.Client
}

func NewRequester(pHost string, pClient *http.Client) IRequester {
	return &sRequester{
		fHost:   pHost,
		fClient: pClient,
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
	if result != hls_settings.CServiceFullName {
		return "", ErrInvalidTitle
	}

	return result, nil
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

func (p *sRequester) FetchRequest(pCtx context.Context, pRequest *hls_settings.SRequest) (response.IResponse, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleNetworkRequestTemplate, p.fHost),
		pRequest,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	resp, err := response.LoadResponse(string(res))
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return resp, nil
}

func (p *sRequester) SendRequest(pCtx context.Context, pRequest *hls_settings.SRequest) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPut,
		fmt.Sprintf(cHandleNetworkRequestTemplate, p.fHost),
		pRequest,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
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

func (p *sRequester) GetOnlines(pCtx context.Context) ([]string, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleNetworkOnlineTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	var onlines []string
	if err := encoding.DeserializeJSON(res, &onlines); err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}

	return onlines, nil
}

func (p *sRequester) DelOnline(pCtx context.Context, pConnect string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleNetworkOnlineTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) GetConnections(pCtx context.Context) ([]string, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	var connects []string
	if err := encoding.DeserializeJSON(res, &connects); err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}

	return connects, nil
}

func (p *sRequester) AddConnection(pCtx context.Context, pConnect string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) DelConnection(pCtx context.Context, pConnect string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		pConnect,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
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
