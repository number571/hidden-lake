package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
	hls_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate           = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleChatMessageTemplate     = "http://" + "%s" + hls_settings.CHandleChatPushPath + "?friend=%s"
	cHandleChatHistorySizeTemplate = "http://" + "%s" + hls_settings.CHandleChatSizePath + "?friend=%s"
	cHandleChatHistoryLoadTemplate = "http://" + "%s" + hls_settings.CHandleChatLoadPath + "?friend=%s&index=%d"
	cHandleChatSubscribeTemplate   = "http://" + "%s" + hls_settings.CHandleChatListenPath + "?friend=%s&sid=%s"
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

func (p *sRequester) PushMessage(pCtx context.Context, pFriend string, pBody string) (time.Time, error) {
	rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleChatMessageTemplate, p.fHost, pFriend),
		pBody,
	)
	if err != nil {
		return time.Time{}, errors.Join(ErrBadRequest, err)
	}
	t, err := dto.ParseTimestamp(string(rsp))
	if err != nil {
		return time.Time{}, errors.Join(ErrDecodeResponse, err)
	}
	return t, nil
}

func (p *sRequester) GetChatSize(pCtx context.Context, pFriend string) (uint64, error) {
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

func (p *sRequester) LoadMessage(pCtx context.Context, pFriend string, pIndex uint64) (dto.IMessage, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChatHistoryLoadTemplate, p.fHost, url.QueryEscape(pFriend), pIndex),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	msg := new(dto.SMessage)
	if err := encoding.DeserializeJSON(res, msg); err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return msg, nil
}

func (p *sRequester) ListenChat(pCtx context.Context, pFriend string, pSid string) (dto.IMessage, error) {
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
		msg := &dto.SMessage{}
		if err := encoding.DeserializeJSON(res, msg); err != nil {
			return nil, errors.Join(ErrDecodeResponse, err)
		}
		return msg, nil
	}
}
