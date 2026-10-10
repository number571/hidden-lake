package database

const (
	errPrefix = "internal/services/notifier/internal/database = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadMessage     = &SError{"load message"}
	ErrExistMessage    = &SError{"exist message"}
	ErrGetMessage      = &SError{"get message"}
	ErrGetIndexMessage = &SError{"get index message"}
	ErrSetMessage      = &SError{"set message"}
	ErrInvalidState    = &SError{"invalid state"}
	ErrSetSizeMessage  = &SError{"set size message"}
	ErrCloseDB         = &SError{"close db"}
	ErrCreateDB        = &SError{"create db"}
	ErrNotFoundMessage = &SError{"not found message"}
)
