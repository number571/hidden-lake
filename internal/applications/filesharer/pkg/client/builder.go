package client

import (
	"fmt"
	"net/http"

	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) GetListFiles(pPage uint64) hls_request.IRequest {
	return hls_request.NewRequest().
		WithMethod(http.MethodGet).
		WithHost(hlf_settings.CServiceFullName).
		WithPath(fmt.Sprintf("%s?page=%d", hlf_settings.CListPath, pPage))
}

func (p *sBuilder) LoadFileChunk(pName string, pChunk uint64) hls_request.IRequest {
	return hls_request.NewRequest().
		WithMethod(http.MethodGet).
		WithHost(hlf_settings.CServiceFullName).
		WithPath(fmt.Sprintf("%s?name=%s&chunk=%d", hlf_settings.CLoadPath, pName, pChunk))
}
