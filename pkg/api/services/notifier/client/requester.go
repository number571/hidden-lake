package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	hls_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/chars"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate           = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleChanPushTemplate        = "http://" + "%s" + hls_settings.CHandleChanMsgboxPushPath + "?friend=%s"
	cHandleChanNotifyTemplate      = "http://" + "%s" + hls_settings.CHandleChanNotifyPath + "?friend=%s"
	cHandleChanPushWithUIDTemplate = "http://" + "%s" + hls_settings.CHandleChanMsgboxPushPath + "?friend=%s&uid=%s"
	cHandleChanLoadRemoteTemplate  = "http://" + "%s" + hls_settings.CHandleChanMsgboxLoadPath + "?friend=%s&uid=%s&remote"
	cHandleChanSizeTemplate        = "http://" + "%s" + hls_settings.CHandleChanMsgboxSizePath + "?friend=%s"
	cHandleChanLoadByEnumTemplate  = "http://" + "%s" + hls_settings.CHandleChanMsgboxLoadPath + "?friend=%s&index=%d"
	cHandleChanLoadByUIDTemplate   = "http://" + "%s" + hls_settings.CHandleChanMsgboxLoadPath + "?friend=%s&uid=%s"
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
	return p.getLimitSize(pCtx, cHandleChanNotifyTemplate)
}

func (p *sRequester) GetMessageLimit(pCtx context.Context) (uint64, error) {
	return p.getLimitSize(pCtx, cHandleChanPushTemplate)
}

func (p *sRequester) getLimitSize(pCtx context.Context, pPathTemplate string) (uint64, error) {
	_, rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(pPathTemplate, p.fHost, "_"),
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

func (p *sRequester) SendNotification(pCtx context.Context, pFriend string, pBody string) (time.Time, error) {
	_, rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleChanNotifyTemplate, p.fHost, pFriend),
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

func (p *sRequester) PushMessageToBox(pCtx context.Context, pFriend string, pUID string, pBody string) (time.Time, error) {
	if !dto.MessageUIDIsValid(pUID) {
		return time.Time{}, ErrSizeMessageUID
	}
	_, rsp, err := api.Request(
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

func (p *sRequester) GetMsgBoxSize(pCtx context.Context, pFriend string) (uint64, error) {
	_, res, err := api.Request(
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

func (p *sRequester) LoadMessageByEnum(pCtx context.Context, pFriend string, pIndex uint64) (dto.IMessage, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChanLoadByEnumTemplate, p.fHost, url.QueryEscape(pFriend), pIndex),
		nil,
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	msg, err := dto.LoadMessage(res)
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return msg, nil
}

func (p *sRequester) LoadMessageByUID(pCtx context.Context, pFriend string, pUID string) (dto.IMessage, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChanLoadByUIDTemplate, p.fHost, url.QueryEscape(pFriend), url.QueryEscape(pUID)),
		nil,
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	msg, err := dto.LoadMessage(res)
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return msg, nil
}

func (p *sRequester) LoadMessageFromBox(pCtx context.Context, pFriend string, pUID string) (string, error) {
	if !dto.MessageUIDIsValid(pUID) {
		return "", ErrSizeMessageUID
	}
	_, rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleChanLoadRemoteTemplate, p.fHost, pFriend, pUID),
		nil,
		nil,
	)
	if err != nil {
		return "", errors.Join(ErrBadRequest, err)
	}
	res := string(rsp)
	if chars.HasNotGraphicCharacters(res) {
		return "", ErrDecodeResponse
	}
	return res, nil
}

func (p *sRequester) ListenChan(pCtx context.Context, pFriend string, pSid string) (string, dto.IMessageBody, error) {
	for {
		headers, res, err := api.Request(
			pCtx,
			p.fClient,
			http.MethodGet,
			fmt.Sprintf(cHandleChanListenTemplate, p.fHost, url.QueryEscape(pFriend), url.QueryEscape(pSid)),
			nil,
			nil,
		)
		if err != nil {
			return "", nil, errors.Join(ErrBadRequest, err)
		}
		if len(res) == 0 {
			continue
		}
		msgBody, err := dto.LoadMessageBody(res)
		if err != nil {
			return "", nil, errors.Join(ErrDecodeResponse, err)
		}
		uid := headers.Get(hls_settings.CHeaderMessageUID)
		return uid, msgBody, nil
	}
}
