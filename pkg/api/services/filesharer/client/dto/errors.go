package dto

const (
	errPrefix = "pkg/api/services/filesharer/client/dto = "
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
