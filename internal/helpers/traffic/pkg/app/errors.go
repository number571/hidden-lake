package app

const (
	errPrefix = "internal/helpers/traffic/pkg/app = "
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
	ErrInitDB  = &SAppError{"init database"}
	ErrClose   = &SAppError{"close"}
)
