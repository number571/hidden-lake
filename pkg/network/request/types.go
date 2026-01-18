package request

import "github.com/number571/go-peer/pkg/types"

type IRequestBuilder interface {
	Build() IRequest

	WithMethod(string) IRequestBuilder
	WithHost(string) IRequestBuilder
	WithPath(string) IRequestBuilder
	WithHead(map[string]string) IRequestBuilder
	WithBody([]byte) IRequestBuilder
}

type IRequest interface {
	types.IConverter

	GetMethod() string
	GetHost() string
	GetPath() string
	GetHead() map[string]string
	GetBody() []byte
}
