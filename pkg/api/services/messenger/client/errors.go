package client

const (
	errPrefix = "pkg/api/services/messenger/client = "
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
