package database

const (
	errPrefix = "internal/database = "
)

type SDatabaseError struct {
	str string
}

func (err *SDatabaseError) Error() string {
	return errPrefix + err.str
}

var (
	ErrGetCount   = &SDatabaseError{"get count"}
	ErrSetCount   = &SDatabaseError{"set count"}
	ErrParseCount = &SDatabaseError{"parse count"}
	ErrGetMessage = &SDatabaseError{"get message"}
	ErrSetMessage = &SDatabaseError{"set message"}
	ErrKeySize    = &SDatabaseError{"key size"}
	ErrMsgSize    = &SDatabaseError{"msg size"}
	ErrDecodeMsg  = &SDatabaseError{"decode msg"}
)
