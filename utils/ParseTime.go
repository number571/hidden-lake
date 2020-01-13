package utils

import (
	"time"
)

func ParseTime(t string) time.Time {
	res, err := time.Parse(time.RFC1123, t)
	if err != nil {
		return time.Unix(0, 0)
	}
	return res
}
