package layer3

const (
	errPrefix = "internal/utils/layer3 = "
)

type SAppError struct {
	str string
}

func (err *SAppError) Error() string {
	return errPrefix + err.str
}

var (
	ErrInvalidBody = &SAppError{"invalid body"}
	ErrInvalidHead = &SAppError{"invalid head"}
)
