package config

const (
	errPrefix = "internal/services/filesharer/pkg/app/config = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrInvalidConfig      = &SError{"invalid config"}
	ErrLoadLogging        = &SError{"load logging"}
	ErrInvalidLogging     = &SError{"invalid logging"}
	ErrInitConfig         = &SError{"init config"}
	ErrDeserializeConfig  = &SError{"deserialize config"}
	ErrReadConfig         = &SError{"read config"}
	ErrConfigNotExist     = &SError{"config not exist"}
	ErrWriteConfig        = &SError{"write config"}
	ErrConfigAlreadyExist = &SError{"config already exist"}
)
