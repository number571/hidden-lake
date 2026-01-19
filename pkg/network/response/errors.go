package response

const (
	errPrefix = "pkg/network/response = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadBytesJoiner = &SError{"load bytes joiner"}
	ErrDecodeResponse  = &SError{"decode response"}
	ErrUnknownType     = &SError{"unknown type"}
)
