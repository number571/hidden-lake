package utils

import (
	"errors"
	"net/url"
)

func GetBoolValueFromQuery(qp url.Values, key string) (bool, error) {
	v, ok := qp[key]
	if !ok {
		// no data available
		return false, nil
	}
	if len(v) == 0 || v[0] == "" {
		// personal exists only as key (without value)
		return true, nil
	}
	switch v[0] {
	case "0", "false":
		return false, nil
	case "1", "true":
		return true, nil
	default:
		return false, errors.New("personal not found") // nolint: err113
	}
}
