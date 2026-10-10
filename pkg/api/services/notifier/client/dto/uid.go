package dto

import "github.com/number571/go-peer/pkg/crypto/random"

const (
	cUIDSize = 22
)

func NewMessageUID() string {
	return random.NewRandom().GetString(cUIDSize)
}

func MessageUIDIsValid(pUID string) bool {
	return len(pUID) == cUIDSize
}
