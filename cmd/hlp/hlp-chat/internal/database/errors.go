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
	ErrGetMessage = &SDatabaseError{"Get message"}
	ErrSetMessage = &SDatabaseError{"set message"}
)
