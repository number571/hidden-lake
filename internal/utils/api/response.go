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
		respBytes   []byte
	)

	switch x := pRes.(type) {
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

	pW.Header().Set("Content-Type", contentType)
	pW.WriteHeader(pRet)

	if _, err := io.Copy(pW, bytes.NewBuffer(respBytes)); err != nil {
		return errors.Join(ErrCopyBytes, err)
	}

	return nil
}
