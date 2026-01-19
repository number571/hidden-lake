package app

const (
	errPrefix = "internal/services/messenger/pkg/app = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning    = &SError{"app running"}
	ErrService    = &SError{"service"}
	ErrInitDB     = &SError{"init database"}
	ErrClose      = &SError{"close"}
	ErrInitConfig = &SError{"init config"}
	ErrSetBuild   = &SError{"set build"}
	ErrMkdirPath  = &SError{"mkdir path"}
)
