package main

const (
	errPrefix = "cmd/hls/hls-filesharer/client/hls-filesharer-cli = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrMkdirPath                 = &SError{"mkdir path"}
	ErrUnknownAction             = &SError{"unknown action"}
	ErrUnknownStorageType        = &SError{"unknown storage type"}
	ErrAvailableOnlyForTypeLocal = &SError{"available only for type local"}
)
