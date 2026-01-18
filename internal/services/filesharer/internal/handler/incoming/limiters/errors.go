package limiters

const (
	errPrefix = "internal/services/filesharer/internal/handler/stream = "
)

type SLimiterError struct {
	str string
}

func (err *SLimiterError) Error() string {
	return errPrefix + err.str
}

var (
	ErrMessageSizeGteLimit = &SLimiterError{"message size >= limit"}
	ErrGetSettingsHLS      = &SLimiterError{"get settings hlk"}
)
