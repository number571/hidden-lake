package utils

import (
	"os"
)

func FileIsExist(filename string) bool {
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return false
    }
    return true
}
