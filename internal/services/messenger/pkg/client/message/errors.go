package message

const (
	errPrefix = "internal/services/messenger/pkg/client/message = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrUnknownMessageType  = &SError{"unknown message type"}
	ErrUnwrapFile          = &SError{"unwrap file"}
	ErrHasNotWritableChars = &SError{"had not writable chars"}
	ErrMessageNull         = &SError{"message null"}
	ErrReadFile            = &SError{"read file"}
	ErrReadFileSize        = &SError{"read file size"}
	ErrGetFormFile         = &SError{"get form file"}
	ErrUploadFile          = &SError{"upload file"}
)
