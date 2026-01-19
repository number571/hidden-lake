package config

const (
	errPrefix = "internal/adapters/http/pkg/app/config = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadLogging        = &SError{"load logging"}
	ErrInvalidLogging     = &SError{"invalid logging"}
	ErrInvalidConfig      = &SError{"invalid config"}
	ErrInitConfig         = &SError{"init config"}
	ErrDeserializeConfig  = &SError{"deserialize config"}
	ErrReadConfig         = &SError{"read config"}
	ErrConfigNotExist     = &SError{"config not exist"}
	ErrWriteConfig        = &SError{"write config"}
	ErrConfigAlreadyExist = &SError{"config already exist"}
	ErrLoadConfig         = &SError{"load config"}
	ErrRebuildConfig      = &SError{"rebuild config"}
	ErrNetworkNotFound    = &SError{"network not found"}
	ErrBuildConfig        = &SError{"build config"}
)
