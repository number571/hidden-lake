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
	cHandleIndexTemplate          = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleRemoteListTemplate     = "http://" + "%s" + hls_settings.CHandleRemoteListPath + "?friend=%s&page=%d&personal=%t"
	cHandleRemoteFileTemplate     = "http://" + "%s" + hls_settings.CHandleRemoteFilePath + "?friend=%s&name=%s&personal=%t"
	cHandleRemoteFileInfoTemplate = "http://" + "%s" + hls_settings.CHandleRemoteFileInfoPath + "?friend=%s&name=%s&personal=%t"
	cHandleLocalListTemplate      = "http://" + "%s" + hls_settings.CHandleLocalListPath + "?friend=%s&page=%d"
	cHandleLocalFileTemplate      = "http://" + "%s" + hls_settings.CHandleLocalFilePath + "?friend=%s&name=%s"
	cHandleLocalFileInfoTemplate  = "http://" + "%s" + hls_settings.CHandleLocalFileInfoPath + "?friend=%s&name=%s"
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

func (p *sRequester) GetRemoteList(pCtx context.Context, pAliasName string, pPage uint64, pPersonal bool) (fileinfo.IFileInfoList, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleRemoteListTemplate, p.fHost, url.QueryEscape(pAliasName), pPage, pPersonal),
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

func (p *sRequester) GetRemoteFile(pW io.Writer, pCtx context.Context, pAliasName string, pFileName string, pPersonal bool) (bool, string, error) {
	headers, err := api.RequestWithWriter(
		pW,
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleRemoteFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName), pPersonal),
		nil,
	)
	if err != nil {
		return false, "", errors.Join(ErrBadRequest, err)
	}
	alreadyDownload := headers.Get(hls_settings.CHeaderInProcess) == hls_settings.CHeaderProcessModeY
	return alreadyDownload, headers.Get(hls_settings.CHeaderFileHash), nil
}

func (p *sRequester) GetRemoteFileInfo(pCtx context.Context, pAliasName string, pFileName string, pPersonal bool) (fileinfo.IFileInfo, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleRemoteFileInfoTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName), pPersonal),
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

func (p *sRequester) GetLocalList(pCtx context.Context, pAliasName string, pPage uint64) (fileinfo.IFileInfoList, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleLocalListTemplate, p.fHost, url.QueryEscape(pAliasName), pPage),
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

func (p *sRequester) GetLocalFile(pW io.Writer, pCtx context.Context, pAliasName string, pFileName string) error {
	_, err := api.RequestWithWriter(
		pW,
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleLocalFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) PutLocalFile(pCtx context.Context, pAliasName string, pFileName string, pR io.Reader) error {
	_, err := api.RequestWithReader(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleLocalFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		pR,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) DelLocalFile(pCtx context.Context, pAliasName string, pFileName string) error {
	_, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleLocalFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sRequester) GetLocalFileInfo(pCtx context.Context, pAliasName string, pFileName string) (fileinfo.IFileInfo, error) {
	res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleLocalFileInfoTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
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
