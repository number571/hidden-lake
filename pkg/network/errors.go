package network

const (
	errPrefix = "pkg/network = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrSendRequest      = &SAppError{"send request"}
	ErrFetchRequest     = &SAppError{"fetch request"}
	ErrLoadResponse     = &SAppError{"load response"}
	ErrAdapterNotRunner = &SAppError{"adapter not runner"}
	ErrRunning          = &SAppError{"node running"}
)
