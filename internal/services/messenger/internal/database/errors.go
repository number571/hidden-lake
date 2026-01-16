package database

const (
	errPrefix = "internal/services/messenger/internal/database = "
)

type SDatabaseError struct {
	str string
}

func (err *SDatabaseError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadMessage    = &SDatabaseError{"load message"}
	ErrGetMessage     = &SDatabaseError{"get message"}
	ErrSetMessage     = &SDatabaseError{"set message"}
	ErrSetSizeMessage = &SDatabaseError{"set size message"}
	ErrCloseDB        = &SDatabaseError{"close db"}
	ErrCreateDB       = &SDatabaseError{"create db"}
)
