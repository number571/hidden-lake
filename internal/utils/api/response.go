package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
)

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

func Response(
	pW http.ResponseWriter,
	pRet int,
	pRes interface{},
) error {
	var (
		contentType string
	)

	switch pRes.(type) {
	case []byte:
		contentType = CApplicationOctetStream
	case string:
		contentType = CTextPlain
	default:
		contentType = CApplicationJSON
	}

	return ResponseWithContentType(pW, pRet, contentType, pRes)
}

func ResponseWithContentType(
	pW http.ResponseWriter,
	pRet int,
	pContentType string,
	pRes interface{},
) error {
	var (
		respBytes []byte
	)

	switch x := pRes.(type) {
	case []byte:
		respBytes = x
	case string:
		respBytes = []byte(x)
	default:
		respBytes = encoding.SerializeJSON(x)
	}

	pW.Header().Set("Content-Type", pContentType)
	pW.WriteHeader(pRet)

	if _, err := io.Copy(pW, bytes.NewBuffer(respBytes)); err != nil {
		return errors.Join(ErrCopyBytes, err)
	}

	return nil
}
