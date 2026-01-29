package database

import (
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/storage/database"
	message "github.com/number571/hidden-lake/pkg/api/services/messenger/client/dto"
)

type sKeyValueDB struct {
	fMutex sync.Mutex
	fDB    database.IKVDatabase
}

func NewKeyValueDB(pPath string) (IKVDatabase, error) {
	db, err := database.NewKVDatabase(pPath)
	if err != nil {
		return nil, errors.Join(ErrCreateDB, err)
	}
	return &sKeyValueDB{fDB: db}, nil
}

func (p *sKeyValueDB) Size(pR IRelation) uint64 {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.getSize(pR)
}

func (p *sKeyValueDB) Load(pR IRelation, pStart, pCount uint64) ([]message.IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	end := min(pStart+pCount, p.getSize(pR))
	res := make([]message.IMessage, 0, pCount)
	for i := pStart; i < end; i++ {
		data, err := p.fDB.Get(getKeyMessageByEnum(pR, i))
		if err != nil {
			return nil, errors.Join(ErrGetMessage, err)
		}
		msg, err := message.LoadMessage(data)
		if err != nil {
			return nil, errors.Join(ErrLoadMessage, err)
		}
		res = append(res, msg)
	}

	return res, nil
}

func (p *sKeyValueDB) Push(pR IRelation, pMsg message.IMessage) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	size := p.getSize(pR)
	numBytes := encoding.Uint64ToBytes(size + 1)
	if err := p.fDB.Set(getKeySize(pR), numBytes[:]); err != nil {
		return errors.Join(ErrSetSizeMessage, err)
	}

	if err := p.fDB.Set(getKeyMessageByEnum(pR, size), pMsg.ToBytes()); err != nil {
		return errors.Join(ErrSetMessage, err)
	}

	return nil
}

func (p *sKeyValueDB) Close() error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if err := p.fDB.Close(); err != nil {
		return errors.Join(ErrCloseDB, err)
	}
	return nil
}

func (p *sKeyValueDB) getSize(pR IRelation) uint64 {
	data, err := p.fDB.Get(getKeySize(pR))
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}
