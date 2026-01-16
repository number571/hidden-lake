package client

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/number571/go-peer/pkg/encoding"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/utils"
	"github.com/number571/hidden-lake/internal/utils/api"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate               = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleStorageListTemplate         = "http://" + "%s" + hls_settings.CHandleStorageListPath + "?friend=%s&page=%d"
	cHandleStorageFileTemplate         = "http://" + "%s" + hls_settings.CHandleStorageFilePath + "?friend=%s&file=%s"
	cHandleStorageFileDonwloadTemplate = "http://" + "%s" + hls_settings.CHandleStorageFilePath + "?friend=%s&file=%s&download"
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

func (p *sRequester) GetFileInfo(pCtx context.Context, pAliasName string, pFileName string) (utils.IFileInfo, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleStorageFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	info := &utils.SFileInfo{}
	if err := encoding.DeserializeJSON(res, info); err != nil {
		return nil, errors.Join(ErrInvalidResponse, err)
	}

	if !isValidHexHash(info.FHash) {
		return nil, ErrInvalidResponse
	}

	return info, nil
}

func (p *sRequester) GetListFiles(pCtx context.Context, pAliasName string, pPage uint64) ([]utils.IFileInfo, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleStorageListTemplate, p.fHost, url.QueryEscape(pAliasName), pPage),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}

	list := make([]utils.SFileInfo, 0, hls_settings.CDefaultPageOffset)
	if err := encoding.DeserializeJSON(res, &list); err != nil {
		return nil, errors.Join(ErrInvalidResponse, err)
	}

	fileInfos := make([]utils.IFileInfo, 0, len(list))
	for _, info := range list {
		if !isValidHexHash(info.FHash) {
			return nil, ErrInvalidResponse
		}
		fileInfos = append(fileInfos, utils.NewFileInfo(info.FName, info.FHash, info.FSize))
	}

	return fileInfos, nil
}

func (p *sRequester) DownloadFile(pW io.Writer, pCtx context.Context, pAliasName string, pFileName string) error {
	err := api.RequestWithWriter(
		pW,
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleStorageFileDonwloadTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func isValidHexHash(hash string) bool {
	v, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	return len(v) == sha512.Size384
}
