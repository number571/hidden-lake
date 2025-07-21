package database

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/storage/database"
	utils_database "github.com/number571/hidden-lake/internal/utils/database"
)

var (
	_ IDatabase = &sVoidDatabase{}
)

type sVoidDatabase struct {
	db database.IKVDatabase
}

func NewVoidDatabase() IDatabase {
	return &sVoidDatabase{
		db: utils_database.NewVoidKVDatabase(),
	}
}

func (p *sVoidDatabase) GetOrigin() database.IKVDatabase                       { return p.db }
func (p *sVoidDatabase) Insert(asymmetric.IPubKey, SMessage) error             { return nil }
func (p *sVoidDatabase) Select(asymmetric.IPubKey, uint64) ([]SMessage, error) { return nil, nil }
