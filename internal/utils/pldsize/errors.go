package pldsize

const (
	errPrefix = "internal/utils/pldsize = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrMessageSizeGteLimit = &SError{"message size >= limit"}
	ErrGetSettingsHLS      = &SError{"get settings hlk"}
)
