package pubkey

const (
	errPrefix = "internal/utils/pubkey = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrGetFriends         = &SError{"get friends"}
	ErrUndefinedPublicKey = &SError{"undefined public key"}
)
