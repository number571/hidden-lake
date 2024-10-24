package adapted

const (
	errPrefix = "internal/adapters/common/internal/producer/internal/adapted = "
)

type SAdaptedError struct {
	str string
}

func (err *SAdaptedError) Error() string {
	return errPrefix + err.str
}

var (
	ErrInvalidResponse = &SAdaptedError{"invalid response"}
	ErrReadResponse    = &SAdaptedError{"read response"}
	ErrBadStatusCode   = &SAdaptedError{"bad status code"}
	ErrBadRequest      = &SAdaptedError{"bad request"}
	ErrBuildRequest    = &SAdaptedError{"build request"}
)
