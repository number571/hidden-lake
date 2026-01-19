package client

const (
	errPrefix = "internal/adapters/http/pkg/client = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest       = &SError{"bad request"}
	ErrDecodeResponse   = &SError{"decode response"}
	ErrInvalidPublicKey = &SError{"invalid public key"}
	ErrInvalidTitle     = &SError{"invalid title"}
)
