package handler

const (
	errPrefix = "pkg/network/handler = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest          = &SError{"bad request"}
	ErrBuildRequest        = &SError{"build request"}
	ErrUndefinedService    = &SError{"undefined service"}
	ErrLoadRequest         = &SError{"load request"}
	ErrInvalidResponseMode = &SError{"invalid response mode"}
)
