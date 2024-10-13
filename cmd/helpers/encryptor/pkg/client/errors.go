package client

const (
	errPrefix = "cmd/helpers/encryptor/pkg/client = "
)

type SClientError struct {
	str string
}

func (err *SClientError) Error() string {
	return errPrefix + err.str
}

var (
	ErrBadRequest       = &SClientError{"bad request"}
	ErrDecodeResponse   = &SClientError{"decode response"}
	ErrInvalidHexFormat = &SClientError{"invalid hex format"}
	ErrInvalidPublicKey = &SClientError{"invalid public key"}
	ErrInvalidTitle     = &SClientError{"invalid title"}
)
