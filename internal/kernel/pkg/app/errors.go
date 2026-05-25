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
	ErrGetScheme        = &SError{"get scheme"}
	ErrInitConfig       = &SError{"init config"}
	ErrSetParallelNull  = &SError{"set parallel = 0"}
	ErrGetParallel      = &SError{"get parallel"}
	ErrGetConsumers     = &SError{"get consumers"}
	ErrCreateAnonNode   = &SError{"create anon node"}
	ErrOpenKVDatabase   = &SError{"open kv database"}
	ErrReadKVDatabase   = &SError{"read kv database"}
	ErrMessageSizeLimit = &SError{"message size limit"}
	ErrSetBuild         = &SError{"set build"}
	ErrMkdirPath        = &SError{"mkdir path"}
	ErrAddFriendToList  = &SError{"add friend to list"}
	ErrCreateNode       = &SError{"create node"}
)
