package dto

import (
	"strings"

	"github.com/number571/go-peer/pkg/crypto/random"
)

const (
	cUIDSize = 22
)

func NewMessageUID() string {
	return random.NewRandom().GetString(cUIDSize)
}

func NewMessageNullUID() string {
	return strings.Repeat("0", cUIDSize)
}

func MessageUIDIsValid(pUID string) bool {
	return len(pUID) == cUIDSize
}
