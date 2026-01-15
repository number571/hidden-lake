package client

import "errors"

const (
	errPrefix = "internal/services/filesharer/pkg/client = "
)

type SClientError struct {
	str string
}

func (err *SClientError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest      = errors.New(errPrefix + "bad request")
	ErrDecodeResponse  = errors.New(errPrefix + "decode response")
	ErrInvalidResponse = errors.New(errPrefix + "invalid response")
	ErrInvalidTitle    = &SClientError{"invalid title"}
)
