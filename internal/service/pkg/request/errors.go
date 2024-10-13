package request

const (
	errPrefix = "internal/service/pkg/request = "
)

type SRequestError struct {
	str string
}

func (err *SRequestError) Error() string {
	return errPrefix + err.str
}

var (
	ErrLoadBytesJoiner = &SRequestError{"load bytes joiner"}
	ErrDecodeRequest   = &SRequestError{"decode request"}
	ErrUnknownType     = &SRequestError{"unknown type"}
)
