package client

const (
	errPrefix = "internal/services/messenger/pkg/client = "
)

type SClientError struct {
	str string
}

func (err *SClientError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest     = &SClientError{"bad request"}
	ErrDecodeResponse = &SClientError{"decode response"}
	ErrPingMessage    = &SClientError{"ping message"}
	ErrPushMessage    = &SClientError{"push message"}
)
