package app

const (
	errPrefix = "internal/applications/remoter/pkg/app = "
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
	ErrClose   = &SAppError{"close"}
)
