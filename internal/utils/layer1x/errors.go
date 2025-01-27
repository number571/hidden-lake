package layer1x

const (
	errPrefix = "internal/utils/layer1x = "
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
