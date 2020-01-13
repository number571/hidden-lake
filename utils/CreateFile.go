package utils

import (
	"errors"
	"os"
)

func CreateFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.New("file not created")
	}
	file.Close()
	return nil
}
