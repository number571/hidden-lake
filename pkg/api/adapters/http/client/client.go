package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/encoding"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/pkg/api/adapters/http/config"
	"github.com/number571/hidden-lake/pkg/network/adapters"
)

var (
	_ IClient = &sClient{}
)

const (
	cHandleIndexTemplate          = "http://" + "%s" + hla_settings.CHandleIndexPath
	cHandleConfigSettingsTemplate = "http://" + "%s" + hla_settings.CHandleConfigSettingsPath
	cHandleConfigConnectsTemplate = "http://" + "%s" + hla_settings.CHandleConfigConnectsPath
	cHandleNetworkOnlineTemplate  = "http://" + "%s" + hla_settings.CHandleNetworkOnlinePath
	cHandleAdapterProduceTemplate = "http://" + "%s" + hla_settings.CHandleAdapterProducePath
	cHandleAdapterConsumeTemplate = "http://" + "%s" + hla_settings.CHandleAdapterConsumePath + "?sid=%s"
)

type sClient struct {
	fHost            string
	fClient          *http.Client
	fAdapterSettings adapters.ISettings
}

func NewClient(pHost string, pClient *http.Client, pAdapterSettings adapters.ISettings) IClient {
	return &sClient{
		fHost:            pHost,
		fClient:          pClient,
		fAdapterSettings: pAdapterSettings,
	}
}

func (p *sClient) GetIndex(pCtx context.Context, pScheme string) error {
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
	index := strings.Index(string(res), "=")
	if index == -1 {
		return ErrInvalidTitle
	}
	if string(res[index+1:]) != pScheme {
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

func (p *sClient) ProduceMessage(pCtx context.Context, pNetMsg layer1.IMessage) error {
	if _, err := layer1.LoadMessage(p.fAdapterSettings, pNetMsg.ToString()); err != nil {
		return errors.Join(ErrDecodeRequest, err)
	}
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleAdapterProduceTemplate, p.fHost),
		nil,
		pNetMsg.ToString(),
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) ConsumeMessage(pCtx context.Context, pSid string) (layer1.IMessage, error) {
	for {
		_, res, err := api.Request(
			pCtx,
			p.fClient,
			http.MethodGet,
			fmt.Sprintf(cHandleAdapterConsumeTemplate, p.fHost, url.QueryEscape(pSid)),
			nil,
			nil,
		)
		if err != nil {
			return nil, errors.Join(ErrBadRequest, err)
		}
		if len(res) == 0 {
			continue
		}
		msg, err := layer1.LoadMessage(p.fAdapterSettings, string(res))
		if err != nil {
			return nil, errors.Join(ErrDecodeResponse, err)
		}
		return msg, nil
	}
}
