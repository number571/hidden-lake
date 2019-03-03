package utils

import (
	"os"
)

func ReadFile(filename string) string {
	file, err := os.Open(filename)
	CheckError(err)
	defer file.Close()

	var (
		buffer []byte = make([]byte, 512)
		data string
	)

	for {
		length, err := file.Read(buffer)
		if length == 0 || err != nil { break }
		data += string(buffer[:length])
	}

	return data
}
