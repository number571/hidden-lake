package utils

import (
	"os"
)

func WriteFile(filename, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	file.WriteString(data)
	file.Close()
	return nil
}
