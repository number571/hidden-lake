package https

const (
	errPrefix = "pkg/network/adapters/https = "
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
	ErrBuildRequest   = &SError{"build request"}
	ErrBadRequest     = &SError{"bad request"}
	ErrDecodeResponse = &SError{"decode response"}
	ErrMessageExist   = &SError{"message exist"}
	ErrNoPassword     = &SError{"no password"}
)
