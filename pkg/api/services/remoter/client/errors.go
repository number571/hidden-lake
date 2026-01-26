package client

const (
	errPrefix = "pkg/api/services/remoter/client = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest     = &SError{"bad request"}
	ErrDecodeResponse = &SError{"decode response"}
	ErrInvalidTitle   = &SError{"invalid title"}
)
