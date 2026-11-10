package main

const (
	errPrefix = "cmd/hls/hls-publisher/client/hls-publisher-cli = "
)

type SError struct {
	str string
}

func (err *SError) Error() string {
	return errPrefix + err.str
}

var (
	ErrUnknownAction                = &SError{"unknown action"}
	ErrUnknownRemoteAlgOption       = &SError{"unknown remote alg option"}
	ErrAvailableOnlyForTypePersonal = &SError{"available only for type personal"}
)
