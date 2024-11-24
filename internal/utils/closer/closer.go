package closer

import (
	"errors"
	"io"
)

// Close all elements in a slice.
func CloseAll(pClosers []io.Closer) error {
	errList := make([]error, 0, len(pClosers))
	for _, c := range pClosers {
		if err := c.Close(); err != nil {
			errList = append(errList, err)
		}
	}
	return errors.Join(errList...)
}
