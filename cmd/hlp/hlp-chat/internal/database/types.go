package database

import (
	"crypto/ed25519"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/storage/database"
)

type IDatabase interface {
	GetOrigin() database.IKVDatabase

	Insert(asymmetric.IPubKey, SMessage) error
	Select(asymmetric.IPubKey, uint64) ([]SMessage, error)
}

type SMessage struct {
	FSendTime time.Time
	FSender   ed25519.PublicKey
	FMessage  string
}
