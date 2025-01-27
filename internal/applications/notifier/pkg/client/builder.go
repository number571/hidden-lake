package client

import (
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	hln_settings "github.com/number571/hidden-lake/internal/applications/notifier/pkg/settings"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

func init() {
	b1 := NewBuilder().Finalyze(0, make([]byte, CSaltSize), []byte{}).ToBytes()
	b2 := NewBuilder().Redirect(0, make([]byte, CSaltSize), []byte{}).ToBytes()
	if len(b1) != len(b2) {
		panic("notifier builder length")
	}
}

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
}

func NewBuilder() IBuilder {
	return &sBuilder{}
}

func (p *sBuilder) Finalyze(pProof uint64, pSalt, pMsg []byte) hls_request.IRequest {
	return hls_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hln_settings.CServiceFullName).
		WithPath(hln_settings.CFinalyzePath).
		WithHead(getHead(pProof, pSalt)).
		WithBody(pMsg).
		Build()
}

func (p *sBuilder) Redirect(pProof uint64, pSalt, pMsg []byte) hls_request.IRequest {
	return hls_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hln_settings.CServiceFullName).
		WithPath(hln_settings.CRedirectPath).
		WithHead(getHead(pProof, pSalt)).
		WithBody(pMsg).
		Build()
}

func getHead(pProof uint64, pSalt []byte) map[string]string {
	proofBytes := encoding.Uint64ToBytes(pProof)
	return map[string]string{
		hln_settings.CHeaderPow:  encoding.HexEncode(proofBytes[:]),
		hln_settings.CHeaderSalt: encoding.HexEncode(pSalt),
	}
}
