package database

const (
	errPrefix = "internal/services/messenger/internal/database = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadMessage    = &SError{"load message"}
	ErrGetMessage     = &SError{"get message"}
	ErrSetMessage     = &SError{"set message"}
	ErrSetSizeMessage = &SError{"set size message"}
	ErrCloseDB        = &SError{"close db"}
	ErrCreateDB       = &SError{"create db"}
)
