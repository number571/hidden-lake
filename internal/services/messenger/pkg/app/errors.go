package app

const (
	errPrefix = "internal/services/messenger/pkg/app = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning    = &SAppError{"app running"}
	ErrService    = &SAppError{"service"}
	ErrInitDB     = &SAppError{"init database"}
	ErrClose      = &SAppError{"close"}
	ErrInitConfig = &SAppError{"init config"}
	ErrSetBuild   = &SAppError{"set build"}
	ErrMkdirPath  = &SAppError{"mkdir path"}
)
