package http

import "net/http"

type sHandlerFunc struct {
	fPath string
	fFunc func(http.ResponseWriter, *http.Request)
}

func NewHandlerFunc(pPath string, pFunc func(http.ResponseWriter, *http.Request)) IHandlerFunc {
	return &sHandlerFunc{
		fPath: pPath,
		fFunc: pFunc,
	}
}

func (p *sHandlerFunc) GetPath() string {
	return p.fPath
}

func (p *sHandlerFunc) GetFunc() func(http.ResponseWriter, *http.Request) {
	return p.fFunc
}
