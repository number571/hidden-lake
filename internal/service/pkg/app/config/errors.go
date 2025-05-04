package config

const (
	errPrefix = "internal/service/pkg/app/config = "
)

type SConfigError struct {
	str string
}

func (err *SConfigError) Error() string {
	return errPrefix + err.str
}

var (
	ErrNotSupportedKeySize = &SConfigError{"not supported key size"}
	ErrInvalidPublicKey    = &SConfigError{"invalid public key"}
	ErrDuplicatePublicKey  = &SConfigError{"duplicate public key"}
	ErrLoadLogging         = &SConfigError{"load logging"}
	ErrInvalidLogging      = &SConfigError{"invalid logging"}
	ErrLoadPublicKey       = &SConfigError{"load public key"}
	ErrInvalidConfig       = &SConfigError{"invalid config"}
	ErrLoadConfig          = &SConfigError{"load config"}
	ErrInitConfig          = &SConfigError{"init config"}
	ErrDeserializeConfig   = &SConfigError{"deserialize config"}
	ErrReadConfig          = &SConfigError{"read config"}
	ErrConfigNotFound      = &SConfigError{"config not found"}
	ErrWriteConfig         = &SConfigError{"write config"}
	ErrConfigAlreadyExist  = &SConfigError{"config already exist"}
	ErrBuildConfig         = &SConfigError{"build config"}
	ErrRebuildConfig       = &SConfigError{"rebuild config"}
	ErrNetworkNotFound     = &SConfigError{"network not found"}
	ErrParseURL            = &SConfigError{"parse url"}
)
