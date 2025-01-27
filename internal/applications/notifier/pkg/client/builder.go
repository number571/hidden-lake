package client

import (
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	hln_settings "github.com/number571/hidden-lake/internal/applications/notifier/pkg/settings"
	hls_request "github.com/number571/hidden-lake/pkg/request"
)

func init() {
	if len(hln_settings.CFinalyzePath) != len(hln_settings.CRedirectPath) {
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
	return basicBuilder(pProof, pSalt, pMsg).
		WithPath(hln_settings.CFinalyzePath).
		Build()
}

func (p *sBuilder) Redirect(pProof uint64, pSalt, pMsg []byte) hls_request.IRequest {
	return basicBuilder(pProof, pSalt, pMsg).
		WithPath(hln_settings.CRedirectPath).
		Build()
}

func basicBuilder(pProof uint64, pSalt, pMsg []byte) hls_request.IRequestBuilder {
	return hls_request.NewRequestBuilder().
		WithMethod(http.MethodPost).
		WithHost(hln_settings.CServiceFullName).
		WithHead(getHead(pProof, pSalt)).
		WithBody(pMsg)
}

func getHead(pProof uint64, pSalt []byte) map[string]string {
	proofBytes := encoding.Uint64ToBytes(pProof)
	return map[string]string{
		hln_settings.CHeaderPow:  encoding.HexEncode(proofBytes[:]),
		hln_settings.CHeaderSalt: encoding.HexEncode(pSalt),
	}
}
