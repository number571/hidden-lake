package fileinfo

const (
	errPrefix = "internal/services/filesharer/pkg/client/fileinfo = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrDecodeInfo  = &SError{"decode info"}
	ErrInvalidHash = &SError{"invalid hash"}
)
