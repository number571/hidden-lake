package utils

import (
	"os"
)

func WriteFile(filename, data string) {
	file, err := os.Create(filename)
	CheckError(err)
	defer file.Close()

	file.WriteString(data)
}
