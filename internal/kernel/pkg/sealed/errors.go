package sealed

const (
	errPrefix = "internal/kernel/pkg/sealed = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrNetworkNotFound         = &SError{"network not found"}
	ErrInvalidParticipantKey   = &SError{"invalid participant key"}
	ErrDuplicateParticipantKey = &SError{"duplicate participant key"}
	ErrLoadLogging             = &SError{"load logging"}
	ErrInvalidSeenKey          = &SError{"invalid seen key"}
	ErrOpenSeenDB              = &SError{"open seen db"}
	ErrDecodeBoot              = &SError{"decode boot"}
	ErrInvalidAPIToken         = &SError{"invalid api token"}
	ErrInvalidPushToken        = &SError{"invalid push token"}
	ErrInvalidSeenDBPath       = &SError{"invalid seen db path"}
	ErrInvalidServiceAddr      = &SError{"invalid service addr"}
	ErrInvalidPrivKey          = &SError{"invalid priv key"}
	ErrUndefinedService        = &SError{"undefined service"}
	ErrUnknownSender           = &SError{"unknown sender"}
	ErrBuildRequest            = &SError{"build request"}
	ErrRequestToService        = &SError{"request to service"}
	ErrInvalidRequest          = &SError{"invalid request (only POST /push)"}
	ErrBoot                    = &SError{"boot"}
	ErrReadBoot                = &SError{"read boot line"}
	ErrCreateNode              = &SError{"create node"}
	ErrListenAPI               = &SError{"listen api"}
	ErrRunNode                 = &SError{"run node"}
	ErrServeAPI                = &SError{"serve api"}
	ErrWriteStdout             = &SError{"write stdout"}
	ErrUnsupportedNetwork      = &SError{"unsupported network (default or not hybrid)"}
	ErrInvalidAliasName        = &SError{"invalid alias name"}
	ErrInvalidEndpoint         = &SError{"invalid endpoint"}
	ErrSeenDBDirNotFound       = &SError{"seen db parent directory not found"}
	ErrSeenDBKeyMismatch       = &SError{"seen db was created with a different key (wrong seen_key or replaced file)"}
)
