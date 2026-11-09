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
	_ IClient = &sClient{}
)

const (
	cHandleIndexTemplate      = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleChanListenTemplate = "http://" + "%s" + hls_settings.CHandleChanListenPath + "?sid=%s"
	cHandleChanNotifyTemplate = "http://" + "%s" + hls_settings.CHandleChanNotifyPath + "?friend=%s&wait=%t"
)

type sClient struct {
	fHost      string
	fClient    *http.Client
	fClientOps *sClientOps
}

func NewClient(pHost string, pClient *http.Client, ops ...IClientOption) IClient {
	clientOps := &sClientOps{}
	for _, o := range ops {
		o(clientOps)
	}
	return &sClient{
		fHost:      pHost,
		fClient:    pClient,
		fClientOps: clientOps,
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
	if string(res) != hls_settings.CAppFullName {
		return ErrInvalidTitle
	}
	return nil
}

func (p *sClient) GetNotifyLimit(pCtx context.Context) (uint64, error) {
	_, rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChanNotifyTemplate, p.fHost, "_", false),
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

func (p *sClient) SendNotification(pCtx context.Context, pFriend string, pBody []byte) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleChanNotifyTemplate, p.fHost, pFriend, p.fClientOps.fWaitResponse),
		nil,
		pBody,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) RecvNotification(pCtx context.Context, pSid string) (string, []byte, error) {
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
