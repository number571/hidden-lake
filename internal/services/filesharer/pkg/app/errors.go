package app

const (
	errPrefix = "internal/services/filesharer/pkg/app = "
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
	ErrInitSTG    = &SAppError{"init storage"}
	ErrClose      = &SAppError{"close"}
	ErrInitConfig = &SAppError{"init config"}
	ErrSetBuild   = &SAppError{"set build"}
	ErrMkdirPath  = &SAppError{"mkdir path"}
)
