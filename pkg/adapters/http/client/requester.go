package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/hidden-lake/internal/utils/api"
	hla_settings "github.com/number571/hidden-lake/pkg/adapters/http/settings"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate          = "http://" + "%s" + hla_settings.CHandleIndexPath
	cHandleConfigSettingsTemplate = "http://" + "%s" + hla_settings.CHandleConfigSettingsPath
	cHandleConfigConnectsTemplate = "http://" + "%s" + hla_settings.CHandleConfigConnectsPath
	cHandleNetworkOnlineTemplate  = "http://" + "%s" + hla_settings.CHandleNetworkOnlinePath
	cHandleNetworkAdapterTemplate = "http://" + "%s" + hla_settings.CHandleNetworkAdapterPath
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
	return string(res), nil
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

func (p *sRequester) ProduceMessage(pCtx context.Context, pNetMsg layer1.IMessage) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleNetworkAdapterTemplate, p.fHost),
		pNetMsg.ToString(),
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}
