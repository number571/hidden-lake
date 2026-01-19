package client

const (
	errPrefix = "internal/services/messenger/pkg/client = "
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
	ErrPingMessage    = &SError{"ping message"}
	ErrPushMessage    = &SError{"push message"}
	ErrInvalidTitle   = &SError{"invalid title"}
)
