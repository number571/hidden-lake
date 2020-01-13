package utils

import (
	"time"
)

func CurrentTime() string {
	return time.Now().Format(time.RFC1123)
}
