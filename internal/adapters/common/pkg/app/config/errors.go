package config

const (
	errPrefix = "internal/adapters/common/consumer/pkg/app/config = "
)

type SConfigError struct {
	str string
}

func (err *SConfigError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadLogging        = &SConfigError{"load logging"}
	ErrInvalidLogging     = &SConfigError{"invalid logging"}
	ErrInvalidConfig      = &SConfigError{"invalid config"}
	ErrInitConfig         = &SConfigError{"init config"}
	ErrDeserializeConfig  = &SConfigError{"deserialize config"}
	ErrReadConfig         = &SConfigError{"read config"}
	ErrConfigNotExist     = &SConfigError{"config not exist"}
	ErrWriteConfig        = &SConfigError{"write config"}
	ErrConfigAlreadyExist = &SConfigError{"config already exist"}
	ErrLoadConfig         = &SConfigError{"load config"}
	ErrRebuildConfig      = &SConfigError{"rebuild config"}
	ErrNetworkNotFound    = &SConfigError{"network not found"}
	ErrBuildConfig        = &SConfigError{"build config"}
)
