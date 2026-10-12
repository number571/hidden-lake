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
	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate           = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleChanPushTemplate        = "http://" + "%s" + hls_settings.CHandleChanPushPath + "?friend=%s"
	cHandleChanPushWithUIDTemplate = "http://" + "%s" + hls_settings.CHandleChanPushPath + "?friend=%s&uid=%s"
	cHandleChanSizeTemplate        = "http://" + "%s" + hls_settings.CHandleChanSizePath + "?friend=%s"
	cHandleChanLoadTemplate        = "http://" + "%s" + hls_settings.CHandleChanLoadPath + "?friend=%s&index=%d"
	cHandleChanListenTemplate      = "http://" + "%s" + hls_settings.CHandleChanListenPath + "?friend=%s&sid=%s"
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
	rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChanPushTemplate, p.fHost, "_"),
		nil,
		nil,
	)
	if err != nil {
		return 0, errors.Join(ErrBadRequest, err)
	}
	var limit [2]uint64
	if err := encoding.DeserializeJSON(rsp, &limit); err != nil {
		return 0, errors.Join(ErrDecodeResponse, err)
	}
	return limit[0], nil
}

func (p *sRequester) GetMessageLimit(pCtx context.Context) (uint64, error) {
	rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChanPushTemplate, p.fHost, "_"),
		nil,
		nil,
	)
	if err != nil {
		return 0, errors.Join(ErrBadRequest, err)
	}
	var limit [2]uint64
	if err := encoding.DeserializeJSON(rsp, &limit); err != nil {
		return 0, errors.Join(ErrDecodeResponse, err)
	}
	return limit[1], nil
}

func (p *sRequester) PushNotify(pCtx context.Context, pFriend string, pBody string) (time.Time, error) {
	rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPut,
		fmt.Sprintf(cHandleChanPushTemplate, p.fHost, pFriend),
		nil,
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

func (p *sRequester) PushMessage(pCtx context.Context, pFriend string, pUID string, pBody string) (time.Time, error) {
	if !dto.MessageUIDIsValid(pUID) {
		return time.Time{}, ErrSizeMessageUID
	}
	rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleChanPushWithUIDTemplate, p.fHost, pFriend, pUID),
		nil,
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
		fmt.Sprintf(cHandleChanSizeTemplate, p.fHost, url.QueryEscape(pFriend)),
		nil,
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
		fmt.Sprintf(cHandleChanLoadTemplate, p.fHost, url.QueryEscape(pFriend), pIndex),
		nil,
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

func (p *sRequester) ListenChan(pCtx context.Context, pFriend string, pSid string) (dto.IMessage, error) {
	for {
		res, err := api.Request(
			pCtx,
			p.fClient,
			http.MethodGet,
			fmt.Sprintf(cHandleChanListenTemplate, p.fHost, url.QueryEscape(pFriend), url.QueryEscape(pSid)),
			nil,
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
