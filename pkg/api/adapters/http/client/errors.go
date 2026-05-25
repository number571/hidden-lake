package client

const (
	errPrefix = "pkg/api/adapters/http/client = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest     = &SError{"bad request"}
	ErrDecodeRequest  = &SError{"decode request"}
	ErrDecodeResponse = &SError{"decode response"}
	ErrInvalidTitle   = &SError{"invalid title"}
)
