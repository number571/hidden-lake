package app

const (
	errPrefix = "internal/applications/filesharer/pkg/app = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning     = &SAppError{"app running"}
	ErrService     = &SAppError{"service"}
	ErrInitSTG     = &SAppError{"init storage"}
	ErrClose       = &SAppError{"close"}
	ErrInitConfig  = &SAppError{"init config"}
	ErrSetNetworks = &SAppError{"set networks"}
)
