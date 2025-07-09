package handler

const (
	errPrefix = "internal/services/filesharer/internal/handler = "
)

type SHandlerError struct {
	str string
}

func (err *SHandlerError) Error() string {
	return errPrefix + err.str
}

var (
	ErrReadOnlineConnections = &SHandlerError{"read online connections"}
	ErrReadConnections       = &SHandlerError{"read connections"}
	ErrGetAllConnections     = &SHandlerError{"get all connections"}
	ErrGetPublicKey          = &SHandlerError{"get public key"}
	ErrGetSettingsHLS        = &SHandlerError{"get settings hlk"}
)
