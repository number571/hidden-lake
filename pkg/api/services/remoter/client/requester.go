package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate       = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleCommandExecTemplate = "http://" + "%s" + hls_settings.CHandleCommandExecPath + "?friend=%s"
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
	if result != hls_settings.CAppFullName {
		return "", ErrInvalidTitle
	}
	return result, nil
}

func (p *sRequester) ExecCommand(pCtx context.Context, pFriend string, pBody *hls_settings.SCommandExecRequest) ([]byte, error) {
	rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleCommandExecTemplate, p.fHost, url.QueryEscape(pFriend)),
		pBody,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	return rsp, nil
}
