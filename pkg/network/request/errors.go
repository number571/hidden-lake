package request

const (
	errPrefix = "pkg/network/request = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadBytesJoiner = &SError{"load bytes joiner"}
	ErrDecodeRequest   = &SError{"decode request"}
	ErrUnknownType     = &SError{"unknown type"}
)
