package http

const (
	errPrefix = "pkg/adapters/http = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning       = &SAppError{"adapter running"}
	ErrNoConnections = &SAppError{"no connections"}
)
