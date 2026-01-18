package fileinfo

const (
	errPrefix = "internal/services/filesharer/pkg/client = "
)

type SFileInfoError struct {
	str string
}

func (err *SFileInfoError) Error() string {
	return errPrefix + err.str
}

var (
	ErrDecodeInfo  = &SFileInfoError{"decode info"}
	ErrInvalidHash = &SFileInfoError{"invalid hash"}
)
