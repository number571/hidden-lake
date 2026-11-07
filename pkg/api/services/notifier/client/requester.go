package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate      = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleChanNotifyTemplate = "http://" + "%s" + hls_settings.CHandleChanNotifyPath + "?friend=%s"
	cHandleChanListenTemplate = "http://" + "%s" + hls_settings.CHandleChanListenPath + "?sid=%s"
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
	if string(res) != hls_settings.CAppFullName {
		return ErrInvalidTitle
	}
	return nil
}

func (p *sRequester) GetNotifyLimit(pCtx context.Context) (uint64, error) {
	_, rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChanNotifyTemplate, p.fHost, "_"),
		nil,
		nil,
	)
	if err != nil {
		return 0, errors.Join(ErrBadRequest, err)
	}
	limit, err := strconv.ParseUint(string(rsp), 10, 64)
	if err != nil {
		return 0, errors.Join(ErrDecodeResponse, err)
	}
	return limit, nil
}

func (p *sRequester) SendNotification(pCtx context.Context, pFriend string, pBody []byte) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleChanNotifyTemplate, p.fHost, pFriend),
		nil,
		pBody,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) RecvNotification(pCtx context.Context, pSid string) (string, []byte, error) {
	for {
		headers, body, err := api.Request(
			pCtx,
			p.fClient,
			http.MethodGet,
			fmt.Sprintf(cHandleChanListenTemplate, p.fHost, url.QueryEscape(pSid)),
			nil,
			nil,
		)
		if err != nil {
			return "", nil, errors.Join(ErrBadRequest, err)
		}
		friend := headers.Get(hls_settings.CHeaderFriendName)
		if friend == "" {
			continue
		}
		return friend, body, nil
	}
}
