package dto

import "time"

func ParseTimestamp(s string) (time.Time, error) {
	return time.Parse(time.DateTime, s)
}
