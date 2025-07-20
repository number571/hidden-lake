package database

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/storage/database"
)

var (
	_ IDatabase            = &sVoidDatabase{}
	_ database.IKVDatabase = &sVoidKVDatabase{}
)

type sVoidDatabase struct {
	db database.IKVDatabase
}
type sVoidKVDatabase struct{}

func NewVoidDatabase() IDatabase {
	return &sVoidDatabase{
		db: newVoidKVDatabase(),
	}
}

func (p *sVoidDatabase) GetOrigin() database.IKVDatabase                       { return p.db }
func (p *sVoidDatabase) Insert(asymmetric.IPubKey, SMessage) error             { return nil }
func (p *sVoidDatabase) Select(asymmetric.IPubKey, uint64) ([]SMessage, error) { return nil, nil }

func newVoidKVDatabase() database.IKVDatabase {
	return &sVoidKVDatabase{}
}

func (p *sVoidKVDatabase) Set([]byte, []byte) error   { return nil }
func (p *sVoidKVDatabase) Get([]byte) ([]byte, error) { return nil, database.ErrNotFound }
func (p *sVoidKVDatabase) Del([]byte) error           { return nil }
func (p *sVoidKVDatabase) Close() error               { return nil }
