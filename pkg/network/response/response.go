package response

import (
	"errors"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/payload/joiner"
)

var (
	_ IResponse = &SResponse{}
)

type SResponse struct {
	SResponseBlock
	FBody []byte `json:"body,omitempty"`
}

type SResponseBlock struct {
	FCode int               `json:"code,omitempty"`
	FHead map[string]string `json:"head,omitempty"`
}

func NewResponseBuilder() IResponseBuilder {
	return &SResponse{
		SResponseBlock: SResponseBlock{},
	}
}

func LoadResponse(pData interface{}) (IResponse, error) {
	var response = new(SResponse)
	switch x := pData.(type) {
	case []byte:
		bytesSlice, err := joiner.LoadBytesJoiner32(x)
		if err != nil || len(bytesSlice) != 2 {
			return nil, ErrLoadBytesJoiner
		}
		if err := encoding.DeserializeJSON(bytesSlice[0], response); err != nil {
			return nil, errors.Join(ErrDecodeResponse, err)
		}
		response.FBody = bytesSlice[1]
	case string:
		if err := encoding.DeserializeJSON([]byte(x), response); err != nil {
			return nil, errors.Join(ErrDecodeResponse, err)
		}
	default:
		return nil, ErrUnknownType
	}
	return response, nil
}

func (p *SResponse) Build() IResponse {
	return p
}

func (p *SResponse) WithCode(pCode int) IResponseBuilder {
	p.FCode = pCode
	return p
}

func (p *SResponse) WithHead(pHead map[string]string) IResponseBuilder {
	p.FHead = make(map[string]string, len(pHead))
	for k, v := range pHead {
		p.FHead[k] = v
	}
	return p
}

func (p *SResponse) WithBody(pBody []byte) IResponseBuilder {
	p.FBody = pBody
	return p
}

func (p *SResponse) ToBytes() []byte {
	return joiner.NewBytesJoiner32([][]byte{
		encoding.SerializeJSON(p.SResponseBlock),
		p.FBody,
	})
}

func (p *SResponse) ToString() string {
	return string(encoding.SerializeJSON(p))
}

func (p *SResponse) GetCode() int {
	return p.FCode
}

func (p *SResponse) GetHead() map[string]string {
	headers := make(map[string]string, len(p.FHead))
	for k, v := range p.FHead {
		headers[k] = v
	}
	return headers
}

func (p *SResponse) GetBody() []byte {
	return p.FBody
}
