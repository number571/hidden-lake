package config

const (
	errPrefix = "internal/kernel/pkg/app/config = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrNotSupportedKeySize     = &SError{"not supported key size"}
	ErrInvalidParticipantKey   = &SError{"invalid participant key"}
	ErrDuplicateParticipantKey = &SError{"duplicate participant key"}
	ErrLoadLogging             = &SError{"load logging"}
	ErrInvalidLogging          = &SError{"invalid logging"}
	ErrLoadParticipantKey      = &SError{"load participant key"}
	ErrInvalidConfig           = &SError{"invalid config"}
	ErrLoadConfig              = &SError{"load config"}
	ErrInitConfig              = &SError{"init config"}
	ErrDeserializeConfig       = &SError{"deserialize config"}
	ErrReadConfig              = &SError{"read config"}
	ErrConfigNotFound          = &SError{"config not found"}
	ErrWriteConfig             = &SError{"write config"}
	ErrConfigAlreadyExist      = &SError{"config already exist"}
	ErrBuildConfig             = &SError{"build config"}
	ErrRebuildConfig           = &SError{"rebuild config"}
	ErrNetworkNotFound         = &SError{"network not found"}
)
