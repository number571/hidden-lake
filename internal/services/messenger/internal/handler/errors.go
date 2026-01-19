package handler

const (
	errPrefix = "internal/services/messenger/internal/handler = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrReadConnections       = &SError{"read connections"}
	ErrReadOnlineConnections = &SError{"read online connections"}
	ErrGetAllConnections     = &SError{"get all connections"}
	ErrGetSettings           = &SError{"get settings"}
	ErrGetPublicKey          = &SError{"get public key"}
	ErrUnknownMessageType    = &SError{"unknown message type"}
	ErrUnwrapFile            = &SError{"unwrap file"}
	ErrHasNotWritableChars   = &SError{"had not writable chars"}
	ErrMessageNull           = &SError{"message null"}
	ErrUndefinedPublicKey    = &SError{"undefined public key"}
	ErrGetFriends            = &SError{"get friends"}
	ErrLenMessageGtLimit     = &SError{"len message > limit"}
	ErrGetMessageLimit       = &SError{"get message limit"}
	ErrPushMessage           = &SError{"push message"}
	ErrReadFile              = &SError{"read file"}
	ErrReadFileSize          = &SError{"read file size"}
	ErrGetFormFile           = &SError{"get form file"}
	ErrUploadFile            = &SError{"upload file"}
)
