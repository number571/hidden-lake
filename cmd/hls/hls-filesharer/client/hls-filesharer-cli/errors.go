package main

const (
	errPrefix = "cmd/hls/hls-filesharer/client/hls-filesharer-cli = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrMkdirPath     = &SError{"mkdir path"}
	ErrRetryNum      = &SError{"retry num"}
	ErrUnknownAction = &SError{"unknown action"}
)
