package std

const (
	errPrefix = "internal/utils/logger/std = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrUnknownLogType = &SError{"unknown log type"}
)
