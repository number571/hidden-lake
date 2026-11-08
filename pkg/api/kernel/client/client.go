package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/encoding"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/pkg/api/kernel/client/scheme"
	"github.com/number571/hidden-lake/pkg/api/kernel/config"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

var (
	_ IClient = &sClient{}
)

const (
	cHandleIndexTemplate          = "http://" + "%s" + hlk_settings.CHandleIndexPath
	cHandleConfigSettingsTemplate = "http://" + "%s" + hlk_settings.CHandleConfigSettingsPath
	cHandleConfigConnectsTemplate = "http://" + "%s" + hlk_settings.CHandleConfigConnectsPath
	cHandleConfigFriendsTemplate  = "http://" + "%s" + hlk_settings.CHandleConfigFriendsPath
	cHandleNetworkOnlineTemplate  = "http://" + "%s" + hlk_settings.CHandleNetworkOnlinePath
	cHandleNetworkRequestTemplate = "http://" + "%s" + hlk_settings.CHandleNetworkRequestPath + "?friend=%s"
)

type sClient struct {
	fHost   string
	fClient *http.Client
}

func NewClient(pHost string, pClient *http.Client) IClient {
	return &sClient{
		fHost:   pHost,
		fClient: pClient,
	}
}

func (p *sClient) GetIndex(pCtx context.Context) error {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	if string(res) != hlk_settings.CAppFullName {
		return ErrInvalidTitle
	}
	return nil
}

func (p *sClient) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigSettingsTemplate, p.fHost),
		nil,
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

func (p *sClient) FetchRequest(pCtx context.Context, pFriend string, pRequest request.IRequest) (response.IResponse, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleNetworkRequestTemplate, p.fHost, url.QueryEscape(pFriend)),
		nil,
		pRequest.ToString(),
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

func (p *sClient) SendRequest(pCtx context.Context, pFriend string, pRequest request.IRequest) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPut,
		fmt.Sprintf(cHandleNetworkRequestTemplate, p.fHost, url.QueryEscape(pFriend)),
		nil,
		pRequest.ToString(),
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) GetFriends(pCtx context.Context, pSchemeType scheme.ISchemeType) (map[string]layer2.IParticipantKey, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		nil,
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	var vFriends [][2]string
	if err := encoding.DeserializeJSON(res, &vFriends); err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}

	result := make(map[string]layer2.IParticipantKey, len(vFriends))
	for _, friend := range vFriends {
		pKey := scheme.LoadParticipantKey(pSchemeType, friend[1])
		if pKey == nil {
			return nil, ErrInvalidParticipantKey
		}
		result[friend[0]] = pKey
	}

	return result, nil
}

func (p *sClient) AddFriend(pCtx context.Context, pFriend string, pKey layer2.IParticipantKey) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		nil,
		[2]string{pFriend, pKey.ToString()},
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) DelFriend(pCtx context.Context, pFriend string) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleConfigFriendsTemplate, p.fHost),
		nil,
		pFriend,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) GetOnlines(pCtx context.Context) ([]string, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleNetworkOnlineTemplate, p.fHost),
		nil,
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

func (p *sClient) DelOnline(pCtx context.Context, pConnect string) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleNetworkOnlineTemplate, p.fHost),
		nil,
		pConnect,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) GetConnections(pCtx context.Context) ([]string, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		nil,
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

func (p *sClient) AddConnection(pCtx context.Context, pConnect string) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		nil,
		pConnect,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) DelConnection(pCtx context.Context, pConnect string) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleConfigConnectsTemplate, p.fHost),
		nil,
		pConnect,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}
