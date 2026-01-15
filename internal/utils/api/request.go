package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
)

func RequestWithWriter(
	pW io.Writer,
	pCtx context.Context,
	pClient *http.Client,
	pMethod, pURL string,
	pData interface{},
) error {
	var (
		contentType string
		reqBytes    []byte
	)

	switch x := pData.(type) {
	case []byte:
		contentType = CTextPlain
		reqBytes = x
	case string:
		contentType = CTextPlain
		reqBytes = []byte(x)
	default:
		contentType = CApplicationJSON
		reqBytes = encoding.SerializeJSON(x)
	}

	req, err := http.NewRequestWithContext(
		pCtx,
		pMethod,
		pURL,
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return errors.Join(ErrBuildRequest, err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := pClient.Do(req)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	defer func() { _ = resp.Body.Close() }()

	_, err = io.Copy(pW, resp.Body)
	return err
}

func Request(
	pCtx context.Context,
	pClient *http.Client,
	pMethod, pURL string,
	pData interface{},
) ([]byte, error) {
	var (
		contentType string
		reqBytes    []byte
	)

	switch x := pData.(type) {
	case []byte:
		contentType = CTextPlain
		reqBytes = x
	case string:
		contentType = CTextPlain
		reqBytes = []byte(x)
	default:
		contentType = CApplicationJSON
		reqBytes = encoding.SerializeJSON(x)
	}

	req, err := http.NewRequestWithContext(
		pCtx,
		pMethod,
		pURL,
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return nil, errors.Join(ErrBuildRequest, err)
	}

	req.Header.Set("Content-Type", contentType)

	resp, err := pClient.Do(req)
	if err != nil {
		return nil, errors.Join(ErrBadRequest, err)
	}
	defer func() { _ = resp.Body.Close() }()

	result, err := loadResponse(resp.StatusCode, resp.Body)
	if err != nil {
		return nil, errors.Join(ErrLoadResponse, err)
	}
	return result, nil
}

func loadResponse(pStatusCode int, pReader io.ReadCloser) ([]byte, error) {
	resp, err := io.ReadAll(pReader)
	if err != nil {
		return nil, errors.Join(ErrReadResponse, err)
	}
	if pStatusCode < 200 || pStatusCode >= 300 {
		return nil, errors.Join(ErrBadStatusCode, fmt.Errorf("status code: %d", pStatusCode)) // nolint: err113
	}
	return resp, nil
}
