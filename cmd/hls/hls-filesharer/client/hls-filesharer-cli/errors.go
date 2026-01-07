package main

const (
	errPrefix = "cmd/hls/hls-filesharer/client/hls-filesharer-cli = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrMkdirPath     = &SAppError{"mkdir path"}
	ErrRetryNum      = &SAppError{"retry num"}
	ErrUnknownAction = &SAppError{"unknown action"}
)
