package app

const (
	errPrefix = "internal/composite/pkg/app = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning        = &SError{"app running"}
	ErrService        = &SError{"service"}
	ErrClose          = &SError{"close"}
	ErrUnknownService = &SError{"unknown service"}
	ErrHasDuplicates  = &SError{"has duplicates"}
	ErrGetRunners     = &SError{"get runners"}
	ErrInitConfig     = &SError{"init config"}
	ErrInitApp        = &SError{"init app"}
	ErrSetBuild       = &SError{"set build"}
	ErrMkdirPath      = &SError{"mkdir path"}
)
