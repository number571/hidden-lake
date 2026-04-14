package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/message/layer1"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/pkg/api/adapters/http/config"
	"github.com/number571/hidden-lake/pkg/network/adapters"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate          = "http://" + "%s" + hla_settings.CHandleIndexPath
	cHandleConfigSettingsTemplate = "http://" + "%s" + hla_settings.CHandleConfigSettingsPath
	cHandleConfigConnectsTemplate = "http://" + "%s" + hla_settings.CHandleConfigConnectsPath
	cHandleNetworkOnlineTemplate  = "http://" + "%s" + hla_settings.CHandleNetworkOnlinePath
	cHandleAdapterProduceTemplate = "http://" + "%s" + hla_settings.CHandleAdapterProducePath
	cHandleAdapterConsumeTemplate = "http://" + "%s" + hla_settings.CHandleAdapterConsumePath + "?sid=%s"
)

type sRequester struct {
	fHost            string
	fClient          *http.Client
	fAdapterSettings adapters.ISettings
}

func NewRequester(pHost string, pClient *http.Client, pAdapterSettings adapters.ISettings) IRequester {
	return &sRequester{
		fHost:            pHost,
		fClient:          pClient,
		fAdapterSettings: pAdapterSettings,
	}
}

func (p *sRequester) GetIndex(pCtx context.Context, pScheme string) error {
	res, err := api.Request(
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

func (p *sRequester) GetSettings(pCtx context.Context) (config.IConfigSettings, error) {
	res, err := api.Request(
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

func (p *sRequester) GetOnlines(pCtx context.Context) ([]string, error) {
	res, err := api.Request(
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

func (p *sRequester) DelOnline(pCtx context.Context, pConnect string) error {
	_, err := api.Request(
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

func (p *sRequester) GetConnections(pCtx context.Context) ([]string, error) {
	res, err := api.Request(
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

func (p *sRequester) AddConnection(pCtx context.Context, pConnect string) error {
	_, err := api.Request(
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

func (p *sRequester) DelConnection(pCtx context.Context, pConnect string) error {
	_, err := api.Request(
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

func (p *sRequester) ProduceMessage(pCtx context.Context, pNetMsg layer1.IMessage) error {
	if _, err := layer1.LoadMessage(p.fAdapterSettings, pNetMsg.ToString()); err != nil {
		return errors.Join(ErrDecodeRequest, err)
	}
	_, err := api.Request(
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

func (p *sRequester) ConsumeMessage(pCtx context.Context, pSid string) (layer1.IMessage, error) {
	for {
		res, err := api.Request(
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
