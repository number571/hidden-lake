package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/utils"
	"github.com/number571/hidden-lake/internal/helpers/loader/pkg/config"
	hll_settings "github.com/number571/hidden-lake/internal/helpers/loader/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
)

const (
	cHandleIndexTemplate           = "%s" + hll_settings.CHandleIndexPath
	cHandleNetworkTransferTemplate = "%s" + hll_settings.CHandleNetworkTransferPath
	cHandleConfigSettingsTemplate  = "%s" + hll_settings.CHandleConfigSettings
)

var (
	_ IRequester = &sRequester{}
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
		return "", utils.MergeErrors(ErrBadRequest, err)
	}

	result := string(res)
	if result != hll_settings.CServiceFullName {
		return "", ErrInvalidTitle
	}

	return result, nil
}

func (p *sRequester) RunTransfer(pCtx context.Context) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleNetworkTransferTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) StopTransfer(pCtx context.Context) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleNetworkTransferTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return utils.MergeErrors(ErrBadRequest, err)
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
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrBadRequest, err)
	}

	cfgSettings := new(config.SConfigSettings)
	if err := encoding.DeserializeJSON(res, cfgSettings); err != nil {
		return nil, utils.MergeErrors(ErrDecodeResponse, err)
	}

	return cfgSettings, nil
}
