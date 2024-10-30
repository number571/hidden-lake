package closer

import (
	"errors"

	"github.com/number571/go-peer/pkg/types"
)

// Close all elements in a slice.
func CloseAll(pClosers []types.ICloser) error {
	errList := make([]error, 0, len(pClosers))
	for _, c := range pClosers {
		if err := c.Close(); err != nil {
			errList = append(errList, err)
		}
	}
	return errors.Join(errList...)
}
