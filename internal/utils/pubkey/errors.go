package pubkey

const (
	errPrefix = "internal/utils/pubkey = "
)

type SUtilsError struct {
	str string
}

func (err *SUtilsError) Error() string {
	return errPrefix + err.str
}

var (
	ErrGetFriends         = &SUtilsError{"get friends"}
	ErrUndefinedPublicKey = &SUtilsError{"undefined public key"}
)
