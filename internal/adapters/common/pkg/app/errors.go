package app

const (
	errPrefix = "internal/adapters/common/pkg/app = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning = &SAppError{"app running"}
	ErrService = &SAppError{"service"}
)