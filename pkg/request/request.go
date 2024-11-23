package request

import (
	"errors"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

var (
	_ IRequest = &SRequest{}
)

type SRequest struct {
	SRequestBlock
	FBody []byte `json:"body,omitempty"`
}

type SRequestBlock struct {
	FMethod string            `json:"method,omitempty"`
	FHost   string            `json:"host,omitempty"`
	FPath   string            `json:"path,omitempty"`
	FHead   map[string]string `json:"head,omitempty"`
}

func NewRequestBuilder() IRequestBuilder {
	return &SRequest{
		SRequestBlock: SRequestBlock{},
	}
}

func LoadRequest(pData interface{}) (IRequest, error) {
	var request = new(SRequest)
	switch x := pData.(type) {
	case []byte:
		bytesSlice, err := joiner.LoadBytesJoiner32(x)
		if err != nil || len(bytesSlice) != 2 {
			return nil, ErrLoadBytesJoiner
		}
		if err := encoding.DeserializeJSON(bytesSlice[0], request); err != nil {
			return nil, errors.Join(ErrDecodeRequest, err)
		}
		request.FBody = bytesSlice[1]
	case string:
		if err := encoding.DeserializeJSON([]byte(x), request); err != nil {
			return nil, errors.Join(ErrDecodeRequest, err)
		}
	default:
		return nil, ErrUnknownType
	}
	return request, nil
}

func (p *SRequest) Build() IRequest {
	return p
}

func (p *SRequest) WithMethod(pMethod string) IRequestBuilder {
	p.FMethod = pMethod
	return p
}

func (p *SRequest) WithHost(pHost string) IRequestBuilder {
	p.FHost = pHost
	return p
}

func (p *SRequest) WithPath(pPath string) IRequestBuilder {
	p.FPath = pPath
	return p
}

func (p *SRequest) WithHead(pHead map[string]string) IRequestBuilder {
	p.FHead = make(map[string]string, len(pHead))
	for k, v := range pHead {
		p.FHead[k] = v
	}
	return p
}

func (p *SRequest) WithBody(pBody []byte) IRequestBuilder {
	p.FBody = pBody
	return p
}

func (p *SRequest) ToBytes() []byte {
	return joiner.NewBytesJoiner32([][]byte{
		encoding.SerializeJSON(p.SRequestBlock),
		p.FBody,
	})
}

func (p *SRequest) ToString() string {
	return string(encoding.SerializeJSON(p))
}

func (p *SRequest) GetHost() string {
	return p.FHost
}

func (p *SRequest) GetPath() string {
	return p.FPath
}

func (p *SRequest) GetMethod() string {
	return p.FMethod
}

func (p *SRequest) GetHead() map[string]string {
	headers := make(map[string]string, len(p.FHead))
	for k, v := range p.FHead {
		headers[k] = v
	}
	return headers
}

func (p *SRequest) GetBody() []byte {
	return p.FBody
}
