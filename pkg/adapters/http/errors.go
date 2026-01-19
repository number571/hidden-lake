package http

const (
	errPrefix = "pkg/adapters/http = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning       = &SError{"adapter running"}
	ErrNoConnections = &SError{"no connections"}
)
