package broker

const (
	errPrefix = "internal/utils/broker = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLimitSubscribers = &SError{"limit subscribers"}
	ErrValutNotFound    = &SError{"value not found"}
	ErrNotRegistered    = &SError{"not registered"}
)
