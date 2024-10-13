package client

const (
	errPrefix = "internal/applications/messenger/pkg/client = "
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
	ErrPushMessage    = &SClientError{"push message"}
)
