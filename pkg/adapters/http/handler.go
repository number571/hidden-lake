package http

import "net/http"

type sHandler struct {
	fPath string
	fFunc func(http.ResponseWriter, *http.Request)
}

func NewHandler(pPath string, pFunc func(http.ResponseWriter, *http.Request)) IHandler {
	return &sHandler{
		fPath: pPath,
		fFunc: pFunc,
	}
}

func (p *sHandler) GetPath() string {
	return p.fPath
}

func (p *sHandler) GetFunc() func(http.ResponseWriter, *http.Request) {
	return p.fFunc
}
