package network

const (
	errPrefix = "pkg/network = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrSendRequest  = &SError{"send request"}
	ErrFetchRequest = &SError{"fetch request"}
	ErrLoadResponse = &SError{"load response"}
	ErrRunning      = &SError{"node running"}
)
