package client

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IRequester = &sRequester{}
)

type sRequester struct {
	fHlkClient hlk_client.IClient
}

func NewRequester(pHlkClient hlk_client.IClient) IRequester {
	return &sRequester{
		fHlkClient: pHlkClient,
	}
}

func (p *sRequester) GetFileInfo(pCtx context.Context, pAliasName string, pRequest hlk_request.IRequest) (IFileInfo, error) {
	resp, err := p.fHlkClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return nil, ErrDecodeResponse
	}

	info := sFileInfo{}
	if err := encoding.DeserializeJSON(resp.GetBody(), &info); err != nil {
		fmt.Println(string(resp.GetBody()))
		return nil, errors.Join(ErrInvalidResponse, err)
	}

	if !isValidHexHash(info.FHash) {
		return nil, ErrInvalidResponse
	}

	return NewFileInfo(info.FName, info.FHash, info.FSize), nil
}

func (p *sRequester) GetListFiles(pCtx context.Context, pAliasName string, pRequest hlk_request.IRequest) ([]IFileInfo, error) {
	resp, err := p.fHlkClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return nil, ErrDecodeResponse
	}

	list := make([]sFileInfo, 0, hls_filesharer_settings.CDefaultPageOffset)
	if err := encoding.DeserializeJSON(resp.GetBody(), &list); err != nil {
		return nil, errors.Join(ErrInvalidResponse, err)
	}

	fileInfos := make([]IFileInfo, 0, len(list))
	for _, info := range list {
		if !isValidHexHash(info.FHash) {
			return nil, ErrInvalidResponse
		}
		fileInfos = append(fileInfos, NewFileInfo(info.FName, info.FHash, info.FSize))
	}

	return fileInfos, nil
}

func (p *sRequester) LoadFileChunk(pCtx context.Context, pAliasName string, pRequest hlk_request.IRequest) ([]byte, error) {
	resp, err := p.fHlkClient.FetchRequest(pCtx, pAliasName, pRequest)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	if resp.GetCode() != http.StatusOK {
		return nil, ErrDecodeResponse
	}

	return resp.GetBody(), nil
}

func isValidHexHash(hash string) bool {
	v, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	return len(v) == sha512.Size384
}
