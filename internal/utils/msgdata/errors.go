package msgdata

const (
	errPrefix = "internal/utils/msgdata = "
)

type SHandlerError struct {
	str string
}

func (err *SHandlerError) Error() string {
	return errPrefix + err.str
}

var (
	ErrUnknownMessageType  = &SHandlerError{"unknown message type"}
	ErrUnwrapFile          = &SHandlerError{"unwrap file"}
	ErrHasNotWritableChars = &SHandlerError{"had not writable chars"}
	ErrMessageNull         = &SHandlerError{"message null"}
	ErrReadFile            = &SHandlerError{"read file"}
	ErrReadFileSize        = &SHandlerError{"read file size"}
	ErrGetFormFile         = &SHandlerError{"get form file"}
	ErrUploadFile          = &SHandlerError{"upload file"}
)
