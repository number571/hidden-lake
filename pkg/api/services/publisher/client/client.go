package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	hls_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
)

var (
	_ IClient = &sClient{}
)

const (
	cHandleIndexTemplate            = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleMsgboxPushTemplate       = "http://" + "%s" + hls_settings.CHandleMsgboxPushPath + "?friend=%s&uid=%s&remote=%t&pub=%t"
	cHandleMsgboxSizeTemplate       = "http://" + "%s" + hls_settings.CHandleMsgboxSizePath + "?friend=%s&remote=%t&pub=%t"
	cHandleMsgboxLoadByEnumTemplate = "http://" + "%s" + hls_settings.CHandleMsgboxLoadPath + "?friend=%s&nid=%d&remote=%t&pub=%t"
	cHandleMsgboxLoadByUIDTemplate  = "http://" + "%s" + hls_settings.CHandleMsgboxLoadPath + "?friend=%s&uid=%s&remote=%t&pub=%t"
	cHandleMsgboxListenTemplate     = "http://" + "%s" + hls_settings.CHandleMsgboxListenPath + "?sid=%s"
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

func (p *sClient) GetPushLimit(pCtx context.Context) (uint64, error) {
	_, rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleMsgboxPushTemplate, p.fHost, "_", "_", false, false),
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

func (p *sClient) GetBoxSize(pCtx context.Context, pFriend string) (uint64, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleMsgboxSizeTemplate, p.fHost, url.QueryEscape(pFriend), p.fClientOps.fRemote, p.fClientOps.fPublic),
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

func (p *sClient) PushMessageToBox(pCtx context.Context, pFriend string, pUID string, pBody []byte) (uint64, time.Time, error) {
	if !dto.MessageUIDIsValid(pUID) {
		return 0, time.Time{}, ErrSizeMessageUID
	}
	if len(pBody) == 0 {
		return 0, time.Time{}, ErrSizeBodyData
	}
	headers, rsp, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleMsgboxPushTemplate, p.fHost, pFriend, pUID, p.fClientOps.fRemote, p.fClientOps.fPublic),
		nil,
		pBody,
	)
	if err != nil {
		return 0, time.Time{}, errors.Join(ErrBadRequest, err)
	}
	t, err := dto.ParseTimestamp(string(rsp))
	if err != nil {
		return 0, time.Time{}, errors.Join(ErrDecodeResponse, err)
	}
	msgNID := headers.Get(hls_settings.CHeaderMessageNID)
	index, err := strconv.ParseUint(msgNID, 10, 64)
	if err != nil {
		return 0, time.Time{}, errors.Join(ErrDecodeResponse, err)
	}
	return index, t, nil
}

func (p *sClient) LoadMessageFromBoxByEnum(pCtx context.Context, pFriend string, pIndex uint64) (dto.IMessage, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleMsgboxLoadByEnumTemplate, p.fHost, url.QueryEscape(pFriend), pIndex, p.fClientOps.fRemote, p.fClientOps.fPublic),
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

func (p *sClient) LoadMessageFromBoxByUID(pCtx context.Context, pFriend string, pUID string) (uint64, dto.IMessage, error) {
	headers, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleMsgboxLoadByUIDTemplate, p.fHost, url.QueryEscape(pFriend), url.QueryEscape(pUID), p.fClientOps.fRemote, p.fClientOps.fPublic),
		nil,
		nil,
	)
	if err != nil {
		return 0, nil, errors.Join(ErrBadRequest, err)
	}
	msg, err := dto.LoadMessage(res)
	if err != nil {
		return 0, nil, errors.Join(ErrDecodeResponse, err)
	}
	if msg.GetHead().GetUID() != pUID {
		return 0, nil, errors.Join(ErrDecodeResponse, err)
	}
	msgNID := headers.Get(hls_settings.CHeaderMessageNID)
	index, err := strconv.ParseUint(msgNID, 10, 64)
	if err != nil {
		return 0, nil, errors.Join(ErrDecodeResponse, err)
	}
	return index, msg, nil
}

func (p *sClient) ListenMessage(pCtx context.Context, pSid string) (string, uint64, dto.IMessage, error) {
	for {
		headers, res, err := api.Request(
			pCtx,
			p.fClient,
			http.MethodGet,
			fmt.Sprintf(cHandleMsgboxListenTemplate, p.fHost, url.QueryEscape(pSid)),
			nil,
			nil,
		)
		if err != nil {
			return "", 0, nil, errors.Join(ErrBadRequest, err)
		}
		if len(res) == 0 {
			continue
		}
		msgBody, err := dto.LoadMessage(res)
		if err != nil {
			return "", 0, nil, errors.Join(ErrDecodeResponse, err)
		}
		friendName := headers.Get(hls_settings.CHeaderFriendName)
		if friendName == "" {
			return "", 0, nil, errors.Join(ErrDecodeResponse, err)
		}
		msgNID := headers.Get(hls_settings.CHeaderMessageNID)
		index, err := strconv.ParseUint(msgNID, 10, 64)
		if err != nil {
			return "", 0, nil, errors.Join(ErrDecodeResponse, err)
		}
		return friendName, index, msgBody, nil
	}
}
