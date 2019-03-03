package utils

import (
	"os"
	"fmt"
)

func PrintError(err string) {
	fmt.Println("[Error]:", err)
	os.Exit(1)
}
