package privkey

const (
	errPrefix = "internal/utils/privkey = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrInvalidPrivateKey = &SAppError{"invalid private key"}
	ErrReadPrivateKey    = &SAppError{"read private key"}
	ErrWritePrivateKey   = &SAppError{"write private key"}
	ErrSizePrivateKey    = &SAppError{"size private key"}
)
