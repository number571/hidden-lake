package stream

const (
	errPrefix = "internal/services/filesharer/internal/handler/stream = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrAppendToTempFile    = &SError{"append to temp file"}
	ErrDeleteTempFile      = &SError{"delete temp file"}
	ErrLoadFileChunk       = &SError{"load file chunk"}
	ErrHashWriteChunk      = &SError{"hash write chunk"}
	ErrInvalidHash         = &SError{"invalid hash"}
	ErrInvalidSize         = &SError{"invalid size"}
	ErrRetryFailed         = &SError{"retry failed"}
	ErrGotAnotherHash      = &SError{"got another hash"}
	ErrInvalidWhence       = &SError{"invalid whence"}
	ErrNegativePosition    = &SError{"negative position"}
	ErrGetMessageLimit     = &SError{"get message limit"}
	ErrGetFileInfo         = &SError{"get file info"}
	ErrReadTempFile        = &SError{"read temp file"}
	ErrInvalidChunkSize    = &SError{"invalid chunk size"}
	ErrInvalidResponseCode = &SError{"invalid response code"}
)
