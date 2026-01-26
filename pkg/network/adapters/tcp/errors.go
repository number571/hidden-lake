package tcp

const (
	errPrefix = "pkg/network/adapters/tcp = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning   = &SError{"adapter running"}
	ErrBroadcast = &SError{"broadcast message"}
)
