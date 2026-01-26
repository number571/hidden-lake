package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

var (
	_ IRequester = &sRequester{}
)

const (
	cHandleIndexTemplate               = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleStorageListTemplate         = "http://" + "%s" + hls_settings.CHandleStorageListPath + "?friend=%s&page=%d"
	cHandleStorageFileInfoTemplate     = "http://" + "%s" + hls_settings.CHandleStorageFileInfoPath + "?friend=%s&name=%s"
	cHandleStorageFileDonwloadTemplate = "http://" + "%s" + hls_settings.CHandleStorageFileDownloadPath + "?friend=%s&name=%s"
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

func (p *sRequester) GetFileInfo(pCtx context.Context, pAliasName string, pFileName string) (fileinfo.IFileInfo, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleStorageFileInfoTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	info, err := fileinfo.LoadFileInfo(res)
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return info, nil
}

func (p *sRequester) GetListFiles(pCtx context.Context, pAliasName string, pPage uint64) ([]fileinfo.IFileInfo, error) {
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
	infos, err := fileinfo.LoadFileInfoList(res)
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return infos, nil
}

func (p *sRequester) DownloadFile(pW io.Writer, pCtx context.Context, pAliasName string, pFileName string) (bool, string, error) {
	headers, err := api.RequestWithWriter(
		pW,
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleStorageFileDonwloadTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
	)
	if err != nil {
		return false, "", errors.Join(ErrBadRequest, err)
	}
	inProcess := headers.Get(hls_settings.CHeaderInProcess) == hls_settings.CHeaderProcessModeY
	return inProcess, headers.Get(hls_settings.CHeaderFileHash), nil
}
