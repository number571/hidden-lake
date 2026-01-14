package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/message"
	hls_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate         = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandlePushMessageTemplate   = "http://" + "%s" + hls_settings.CHandlePushMessagePath + "?friend=%s"
	cHandleLoadMessagesTemplate  = "http://" + "%s" + hls_settings.CHandleLoadMessagesPath + "?friend=%s&page=%d&offset=%d"
	cHandleListenMessageTemplate = "http://" + "%s" + hls_settings.CHandleListenMessagePath + "?friend=%s&sid=%s"
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

func (p *sRequester) PushMessage(pCtx context.Context, pFriend string, pBody string) (string, error) {
	rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandlePushMessageTemplate, p.fHost, pFriend),
		pBody,
	)
	if err != nil {
		return "", errors.Join(ErrBadRequest, err)
	}
	return string(rsp), nil
}

func (p *sRequester) LoadMessages(pCtx context.Context, pFriend string, pPage uint64, pOffset uint64) ([]message.IMessage, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleLoadMessagesTemplate, p.fHost, url.QueryEscape(pFriend), pPage, pOffset),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	var msgs []*message.SMessage
	if err := encoding.DeserializeJSON(res, &msgs); err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	result := make([]message.IMessage, 0, len(msgs))
	for _, m := range msgs {
		result = append(result, m)
	}
	return result, nil
}

func (p *sRequester) ListenMessage(pCtx context.Context, pFriend string, pSid string) (message.IMessage, error) {
	for {
		res, err := api.Request(
			pCtx,
			p.fClient,
			http.MethodGet,
			fmt.Sprintf(cHandleListenMessageTemplate, p.fHost, url.QueryEscape(pFriend), url.QueryEscape(pSid)),
			nil,
		)
		if err != nil {
			return nil, errors.Join(ErrBadRequest, err)
		}
		if len(res) == 0 {
			continue
		}
		msg := &message.SMessage{}
		if err := encoding.DeserializeJSON(res, msg); err != nil {
			return nil, errors.Join(ErrDecodeResponse, err)
		}
		return msg, nil
	}
}
