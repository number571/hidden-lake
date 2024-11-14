package handler

const (
	errPrefix = "pkg/handler = "
)

type SHandlerError struct {
	str string
}

func (err *SHandlerError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest          = &SHandlerError{"bad request"}
	ErrBuildRequest        = &SHandlerError{"build request"}
	ErrUndefinedService    = &SHandlerError{"undefined service"}
	ErrLoadRequest         = &SHandlerError{"load request"}
	ErrInvalidResponseMode = &SHandlerError{"invalid response mode"}
)
