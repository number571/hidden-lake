package http

const (
	errPrefix = "pkg/network/adapters/http = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning        = &SError{"adapter running"}
	ErrNoConnections  = &SError{"no connections"}
	ErrBadRequest     = &SError{"bad request"}
	ErrDecodeResponse = &SError{"decode response"}
	ErrMessageExist   = &SError{"message exist"}
)
