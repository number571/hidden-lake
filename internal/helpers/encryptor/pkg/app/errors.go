package app

const (
	errPrefix = "internal/helpers/encryptor/pkg/app = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning          = &SAppError{"app running"}
	ErrService          = &SAppError{"service"}
	ErrClose            = &SAppError{"close"}
	ErrSizePrivateKey   = &SAppError{"size private key"}
	ErrGetPrivateKey    = &SAppError{"get private key"}
	ErrInitConfig       = &SAppError{"init config"}
	ErrSetParallelNull  = &SAppError{"set parallel = 0"}
	ErrGetParallelValue = &SAppError{"get parallel value"}
)
