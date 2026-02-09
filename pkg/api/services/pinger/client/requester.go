package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate       = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleCommandPingTemplate = "http://" + "%s" + hls_settings.CHandleCommandPingPath + "?friend=%s"
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

func (p *sRequester) GetIndex(pCtx context.Context) error {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleIndexTemplate, p.fHost),
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	if string(res) != hls_settings.CAppFullName {
		return ErrInvalidTitle
	}
	return nil
}

func (p *sRequester) PingFriend(pCtx context.Context, pFriend string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleCommandPingTemplate, p.fHost, url.QueryEscape(pFriend)),
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}
