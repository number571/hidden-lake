package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/client/message"
	hls_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate           = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleChatMessageTemplate     = "http://" + "%s" + hls_settings.CHandleChatMessagePath + "?friend=%s"
	cHandleChatHistorySizeTemplate = "http://" + "%s" + hls_settings.CHandleChatHistorySizePath + "?friend=%s"
	cHandleChatHistoryLoadTemplate = "http://" + "%s" + hls_settings.CHandleChatHistoryLoadPath + "?friend=%s&start=%d&count=%d&select=%s"
	cHandleChatSubscribeTemplate   = "http://" + "%s" + hls_settings.CHandleChatSubscribePath + "?friend=%s&sid=%s"
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

func (p *sRequester) GetMessageLimit(pCtx context.Context) (uint64, error) {
	rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChatMessageTemplate, p.fHost, "_"),
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

func (p *sRequester) PushMessage(pCtx context.Context, pFriend string, pBody string) (string, error) {
	rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleChatMessageTemplate, p.fHost, pFriend),
		pBody,
	)
	if err != nil {
		return "", errors.Join(ErrBadRequest, err)
	}
	return string(rsp), nil
}

func (p *sRequester) CountMessages(pCtx context.Context, pFriend string) (uint64, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChatHistorySizeTemplate, p.fHost, url.QueryEscape(pFriend)),
		nil,
	)
	if err != nil {
		return 0, errors.Join(ErrBadRequest, err)
	}
	count, err := strconv.ParseUint(string(res), 10, 64)
	if err != nil {
		return 0, errors.Join(ErrDecodeResponse, err)
	}
	return count, nil
}

func (p *sRequester) LoadMessages(pCtx context.Context, pFriend string, pStart uint64, pCount uint64, pDesc bool) ([]message.IMessage, error) {
	selectType := "asc"
	if pDesc {
		selectType = "desc"
	}
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChatHistoryLoadTemplate, p.fHost, url.QueryEscape(pFriend), pStart, pCount, selectType),
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

func (p *sRequester) ListenChat(pCtx context.Context, pFriend string, pSid string) (message.IMessage, error) {
	for {
		res, err := api.Request(
			pCtx,
			p.fClient,
			http.MethodGet,
			fmt.Sprintf(cHandleChatSubscribeTemplate, p.fHost, url.QueryEscape(pFriend), url.QueryEscape(pSid)),
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
