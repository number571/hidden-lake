package client

import (
	"fmt"
	"net/http"

	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hlk_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) GetFileInfo(pFileName string) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_filesharer_settings.CAppShortName).
		WithPath(fmt.Sprintf("%s?name=%s", hls_filesharer_settings.CInfoPath, pFileName)).
		Build()
}

func (p *sBuilder) GetListFiles(pPage uint64) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_filesharer_settings.CAppShortName).
		WithPath(fmt.Sprintf("%s?page=%d", hls_filesharer_settings.CListPath, pPage)).
		Build()
}

func (p *sBuilder) LoadFileChunk(pName string, pChunk uint64) hlk_request.IRequest {
	return hlk_request.NewRequestBuilder().
		WithMethod(http.MethodGet).
		WithHost(hls_filesharer_settings.CAppShortName).
		WithPath(fmt.Sprintf("%s?name=%s&chunk=%d", hls_filesharer_settings.CLoadPath, pName, pChunk)).
		Build()
}
