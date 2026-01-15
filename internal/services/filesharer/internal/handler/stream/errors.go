package stream

const (
	errPrefix = "internal/services/filesharer/internal/handler/stream = "
)

type SStreamError struct {
	str string
}

func (err *SStreamError) Error() string {
	return errPrefix + err.str
}

var (
	ErrAppendToTempFile = &SStreamError{"append to temp file"}
	ErrDeleteTempFile   = &SStreamError{"delete temp file"}
	ErrLoadFileChunk    = &SStreamError{"load file chunk"}
	ErrHashWriteChunk   = &SStreamError{"hash write chunk"}
	ErrInvalidHash      = &SStreamError{"invalid hash"}
	ErrRetryFailed      = &SStreamError{"retry failed"}
	ErrInvalidWhence    = &SStreamError{"invalid whence"}
	ErrNegativePosition = &SStreamError{"negative position"}
	ErrGetMessageLimit  = &SStreamError{"get message limit"}
	ErrGetFileInfo      = &SStreamError{"get file info"}
	ErrReadTempFile     = &SStreamError{"read temp file"}
)
