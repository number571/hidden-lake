package tcp

const (
	errPrefix = "pkg/adapters/tcp = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning   = &SAppError{"adapter running"}
	ErrBroadcast = &SAppError{"broadcast message"}
)
