package dto

const (
	errPrefix = "pkg/api/services/notifier/client/dto = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrUnknownType = &SError{"unknown type"}
)
