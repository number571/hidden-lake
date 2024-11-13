package client

import (
	"context"
	"crypto/sha512"
	"errors"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHLSClient hls_client.IClient
}

func NewRequester(pHLSClient hls_client.IClient) IRequester {
	return &sRequester{
		fHLSClient: pHLSClient,
	}
}

func (p *sRequester) GetListFiles(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) ([]hlf_settings.SFileInfo, error) {
	resp, err := p.fHLSClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return nil, ErrDecodeResponse
	}

	list := make([]hlf_settings.SFileInfo, 0, hlf_settings.CDefaultPageOffset)
	if err := encoding.DeserializeJSON(resp.GetBody(), &list); err != nil {
		return nil, errors.Join(ErrInvalidResponse, err)
	}

	for _, info := range list {
		if len(info.FHash) != (sha512.Size384 << 1) {
			return nil, ErrInvalidResponse
		}
	}

	return list, nil
}

func (p *sRequester) LoadFileChunk(pCtx context.Context, pAliasName string, pRequest hls_request.IRequest) ([]byte, error) {
	resp, err := p.fHLSClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return nil, ErrDecodeResponse
	}

	return resp.GetBody(), nil
}
