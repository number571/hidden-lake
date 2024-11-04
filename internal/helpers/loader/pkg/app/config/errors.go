package config

const (
	errPrefix = "internal/helpers/loader/pkg/app/config = "
)

type SConfigError struct {
	str string
}

func (err *SConfigError) Error() string {
	return errPrefix + err.str
}

var (
	ErrInvalidConfig      = &SConfigError{"invalid config"}
	ErrLoadLogging        = &SConfigError{"load logging"}
	ErrInvalidLogging     = &SConfigError{"invalid logging"}
	ErrInitConfig         = &SConfigError{"init config"}
	ErrDeserializeConfig  = &SConfigError{"deserialize config"}
	ErrReadConfig         = &SConfigError{"read config"}
	ErrConfigNotExist     = &SConfigError{"config not exist"}
	ErrWriteConfig        = &SConfigError{"write config"}
	ErrConfigAlreadyExist = &SConfigError{"config already exist"}
	ErrRebuildConfig      = &SConfigError{"rebuild config"}
	ErrNetworkNotFound    = &SConfigError{"network not found"}
	ErrLoadConfig         = &SConfigError{"load config"}
	ErrBuildConfig        = &SConfigError{"build config"}
)