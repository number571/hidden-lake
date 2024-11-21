package response

import "github.com/number571/go-peer/pkg/types"

type IResponseBuilder interface {
	Build() IResponse

	WithCode(pCode int) IResponseBuilder
	WithHead(map[string]string) IResponseBuilder
	WithBody(pBody []byte) IResponseBuilder
}

type IResponse interface {
	types.IConverter

	GetCode() int
	GetHead() map[string]string
	GetBody() []byte
}
