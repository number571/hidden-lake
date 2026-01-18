package client

const (
	errPrefix = "internal/services/filesharer/pkg/client = "
)

type SClientError struct {
	str string
}

func (err *SClientError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest      = &SClientError{"bad request"}
	ErrDecodeResponse  = &SClientError{"decode response"}
	ErrInvalidResponse = &SClientError{"invalid response"}
	ErrInvalidTitle    = &SClientError{"invalid title"}
)
