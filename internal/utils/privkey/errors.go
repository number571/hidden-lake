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
	ErrInvalidPrivateKey           = &SError{"invalid private key"}
	ErrInvalidPublicKey            = &SError{"invalid public key"}
	ErrReadPrivateKey              = &SError{"read private key"}
	ErrReadPublicKey               = &SError{"read public key"}
	ErrWritePrivateKey             = &SError{"write private key"}
	ErrWritePublicKey              = &SError{"write public key"}
	ErrSizePrivateKey              = &SError{"size private key"}
	ErrNotLinkedPublicKeyToPrivate = &SError{"not linked public key to private"}
)
