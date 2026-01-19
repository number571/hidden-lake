package api

const (
	errPrefix = "internal/utils/api = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadStatusCode = &SError{"bad status code"}
	ErrReadResponse  = &SError{"read response"}
	ErrLoadResponse  = &SError{"load response"}
	ErrBadRequest    = &SError{"bad request"}
	ErrBuildRequest  = &SError{"build request"}
	ErrCopyBytes     = &SError{"copy bytes"}
)
