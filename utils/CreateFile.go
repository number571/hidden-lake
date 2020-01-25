package utils

import (
	"os"
)

func CreateFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}
