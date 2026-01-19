package app

const (
	errPrefix = "internal/kernel/pkg/app = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrRunning          = &SError{"app running"}
	ErrService          = &SError{"service"}
	ErrInitDB           = &SError{"init database"}
	ErrClose            = &SError{"close"}
	ErrSizePrivateKey   = &SError{"size private key"}
	ErrGetPrivateKey    = &SError{"get private key"}
	ErrInitConfig       = &SError{"init config"}
	ErrSetParallelNull  = &SError{"set parallel = 0"}
	ErrGetParallel      = &SError{"get parallel"}
	ErrGetConsumers     = &SError{"get consumers"}
	ErrCreateAnonNode   = &SError{"create anon node"}
	ErrOpenKVDatabase   = &SError{"open kv database"}
	ErrReadKVDatabase   = &SError{"read kv database"}
	ErrMessageSizeLimit = &SError{"message size limit"}
	ErrInvalidPsdPubKey = &SError{"invalid psd public key"}
	ErrGetPsdPubKey     = &SError{"get psd pub key"}
	ErrSetPsdPubKey     = &SError{"set psd pub key"}
	ErrSetBuild         = &SError{"set build"}
	ErrMkdirPath        = &SError{"mkdir path"}
)
