package dto

import (
	"strconv"
)

func ParseRemoteType(v string) (IRemoteType, error) {
	if v == "" {
		return CRemoteDisabled, nil
	}
	t, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0, err
	}
	if t >= uint64(cEndRemoteUndefined) {
		return 0, ErrUnknownType
	}
	return IRemoteType(t), nil
}
