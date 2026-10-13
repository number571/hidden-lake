package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	cContentTypeHeader = "Content-Type"
)

func Response(
	pW http.ResponseWriter,
	pRet int,
	pRsp interface{},
	pCalls ...func(http.ResponseWriter),
) error {
	var (
		contentType string
		respBytes   []byte
	)

	switch x := pRsp.(type) {
	case []byte:
		contentType = CApplicationOctetStream
		respBytes = x
	case string:
		contentType = CTextPlain
		respBytes = []byte(x)
	default:
		contentType = CApplicationJSON
		respBytes = encoding.SerializeJSON(x)
	}

	pW.Header().Set(cContentTypeHeader, contentType)

	for _, call := range pCalls {
		call(pW)
	}

	return ResponseWithReader(pW, pRet, bytes.NewBuffer(respBytes))
}

func ResponseWithReader(
	pW http.ResponseWriter,
	pRet int,
	pR io.Reader,
) error {
	pW.WriteHeader(pRet)

	if _, err := io.Copy(pW, pR); err != nil {
		return errors.Join(ErrCopyBytes, err)
	}

	return nil
}

func WithoutContentType() func(http.ResponseWriter) {
	return func(pW http.ResponseWriter) {
		// INFO: if use Del() func -> Go autoset Content-Type
		// Content-Type full deleted on HLK response side (if value = "")
		// => work only in the HLS applications
		pW.Header().Set(cContentTypeHeader, "")
	}
}

func WithCustomContentType(pContentType string) func(http.ResponseWriter) {
	return func(pW http.ResponseWriter) {
		pW.Header().Set(cContentTypeHeader, pContentType)
	}
}
