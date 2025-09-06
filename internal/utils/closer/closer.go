package closer

import (
	"errors"
	"io"
)

type sCloser struct {
	fClosers []io.Closer
}

func NewCloser(pClosers ...io.Closer) io.Closer {
	return &sCloser{fClosers: pClosers}
}

// Close all elements in a slice.
func (p *sCloser) Close() error {
	errList := make([]error, 0, len(p.fClosers))
	for _, c := range p.fClosers {
		if err := c.Close(); err != nil {
			errList = append(errList, err)
		}
	}
	return errors.Join(errList...)
}
