package http

import (
	"net/http"
)

type sLogBuilder struct {
	fService string
	fMethod  string
	fPath    string
	fConn    string
	fMessage string
}

func NewLogBuilder(pService string, pR *http.Request) ILogBuilder {
	return &sLogBuilder{
		fService: pService,
		fMethod:  pR.Method,
		fPath:    pR.URL.Path,
		fConn:    pR.RemoteAddr,
	}
}

func (p *sLogBuilder) Build() ILogGetter {
	return p
}

func (p *sLogBuilder) WithMessage(pMsg string) ILogBuilder {
	p.fMessage = pMsg
	return p
}

func (p *sLogBuilder) GetService() string {
	return p.fService
}

func (p *sLogBuilder) GetConn() string {
	return p.fConn
}

func (p *sLogBuilder) GetMethod() string {
	return p.fMethod
}

func (p *sLogBuilder) GetPath() string {
	return p.fPath
}

func (p *sLogBuilder) GetMessage() string {
	return p.fMessage
}
