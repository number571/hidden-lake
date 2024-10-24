package config

const (
	errPrefix = "internal/helpers/encryptor/internal/config = "
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
	ErrDuplicatePublicKey = &SConfigError{"duplicate public key"}
	ErrInvalidPublicKey   = &SConfigError{"invalid public key"}
	ErrLoadPublicKey      = &SConfigError{"load public key"}
	ErrLoadConfig         = &SConfigError{"load config"}
)
