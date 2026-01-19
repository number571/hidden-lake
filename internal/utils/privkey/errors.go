package privkey

const (
	errPrefix = "internal/utils/privkey = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrInvalidPrivateKey = &SError{"invalid private key"}
	ErrReadPrivateKey    = &SError{"read private key"}
	ErrWritePrivateKey   = &SError{"write private key"}
	ErrSizePrivateKey    = &SError{"size private key"}
)
