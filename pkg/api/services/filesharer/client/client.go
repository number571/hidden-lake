package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/number571/go-peer/pkg/encoding"
	hls_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

var (
	_ IClient = &sClient{}
)

const (
	cHandleIndexTemplate          = "http://" + "%s" + hls_settings.CHandleIndexPath
	cHandleRemoteListTemplate     = "http://" + "%s" + hls_settings.CHandleRemoteListPath + "?friend=%s&page=%d&pub=%t"
	cHandleRemoteFileTemplate     = "http://" + "%s" + hls_settings.CHandleRemoteFilePath + "?friend=%s&name=%s&pub=%t"
	cHandleRemoteFileInfoTemplate = "http://" + "%s" + hls_settings.CHandleRemoteFileInfoPath + "?friend=%s&name=%s&pub=%t"
	cHandleRemoteFileProcTemplate = "http://" + "%s" + hls_settings.CHandleRemoteFileProcPath + "?friend=%s&name=%s&pub=%t"
	cHandleRemoteListProcTemplate = "http://" + "%s" + hls_settings.CHandleRemoteListProcPath
	cHandleLocalListTemplate      = "http://" + "%s" + hls_settings.CHandleLocalListPath + "?friend=%s&page=%d"
	cHandleLocalFileTemplate      = "http://" + "%s" + hls_settings.CHandleLocalFilePath + "?friend=%s&name=%s"
	cHandleLocalFileInfoTemplate  = "http://" + "%s" + hls_settings.CHandleLocalFileInfoPath + "?friend=%s&name=%s"
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

func (p *sClient) GetRemoteList(pCtx context.Context, pAliasName string, pPage uint64) (uint64, []string, error) {
	header, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleRemoteListTemplate, p.fHost, url.QueryEscape(pAliasName), pPage, p.fClientOps.fPublic),
		nil,
		nil,
	)
	if err != nil {
		return 0, nil, errors.Join(ErrBadRequest, err)
	}
	var list []string
	if err := encoding.DeserializeJSON(res, &list); err != nil {
		return 0, nil, errors.Join(ErrDecodeResponse, err)
	}
	pagesCount, err := strconv.ParseUint(header.Get(hls_settings.CHeaderPagesCount), 10, 64)
	if err != nil {
		return 0, nil, errors.Join(ErrDecodeResponse, err)
	}
	return pagesCount, list, nil
}

func (p *sClient) GetRemoteFile(pW io.Writer, pCtx context.Context, pAliasName string, pFileName string) (bool, error) {
	headers, err := api.RequestWithWriter(
		pW,
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleRemoteFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName), p.fClientOps.fPublic),
		nil,
		nil,
	)
	if err != nil {
		return false, errors.Join(ErrBadRequest, err)
	}
	inProcess := headers.Get(hls_settings.CHeaderInProcess) == hls_settings.CHeaderProcessModeY
	return inProcess, nil
}

func (p *sClient) DelRemoteFile(pCtx context.Context, pAliasName string, pFileName string) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleRemoteFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName), p.fClientOps.fPublic),
		nil,
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) GetRemoteFileInfo(pCtx context.Context, pAliasName string, pFileName string) (dto.IFileInfo, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleRemoteFileInfoTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName), p.fClientOps.fPublic),
		nil,
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	info, err := dto.LoadFileInfo(res)
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return info, nil
}

func (p *sClient) GetLocalList(pCtx context.Context, pAliasName string, pPage uint64) (uint64, []string, error) {
	header, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleLocalListTemplate, p.fHost, url.QueryEscape(pAliasName), pPage),
		nil,
		nil,
	)
	if err != nil {
		return 0, nil, errors.Join(ErrBadRequest, err)
	}
	var list []string
	if err := encoding.DeserializeJSON(res, &list); err != nil {
		return 0, nil, errors.Join(ErrDecodeResponse, err)
	}
	pagesCount, err := strconv.ParseUint(header.Get(hls_settings.CHeaderPagesCount), 10, 64)
	if err != nil {
		return 0, nil, errors.Join(ErrDecodeResponse, err)
	}
	return pagesCount, list, nil
}

func (p *sClient) GetLocalFile(pW io.Writer, pCtx context.Context, pAliasName string, pFileName string) error {
	_, err := api.RequestWithWriter(
		pW,
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleLocalFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) PutLocalFile(pCtx context.Context, pAliasName string, pFileName string, pR io.Reader) error {
	_, _, err := api.RequestWithReader(
		pCtx,
		p.fClient,
		http.MethodPost,
		fmt.Sprintf(cHandleLocalFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
		pR,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) DelLocalFile(pCtx context.Context, pAliasName string, pFileName string) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleLocalFileTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) GetLocalFileInfo(pCtx context.Context, pAliasName string, pFileName string) (dto.IFileInfo, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleLocalFileInfoTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName)),
		nil,
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	info, err := dto.LoadFileInfo(res)
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return info, nil
}

func (p *sClient) GetRemoteFileProc(pCtx context.Context, pAliasName string, pFileName string) (dto.IDownloadProcess, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleRemoteFileProcTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName), p.fClientOps.fPublic),
		nil,
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	dp, err := dto.LoadDownloadProcess(res)
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return dp, nil
}

func (p *sClient) DelRemoteFileProc(pCtx context.Context, pAliasName string, pFileName string) error {
	_, _, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodDelete,
		fmt.Sprintf(cHandleRemoteFileProcTemplate, p.fHost, url.QueryEscape(pAliasName), url.QueryEscape(pFileName), p.fClientOps.fPublic),
		nil,
		nil,
	)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	return nil
}

func (p *sClient) GetRemoteListProc(pCtx context.Context) ([]dto.IDownloadProcess, error) {
	_, res, err := api.Request(
		pCtx,
		p.fClient,
		http.MethodGet,
		fmt.Sprintf(cHandleRemoteListProcTemplate, p.fHost),
		nil,
		nil,
	)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	dpl, err := dto.LoadDownloadProcessList(res)
	if err != nil {
		return nil, errors.Join(ErrDecodeResponse, err)
	}
	return dpl, nil
}
