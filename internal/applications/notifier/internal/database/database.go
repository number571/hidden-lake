package database

import (
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/storage/database"
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

func (p *sKeyValueDB) SetHash(pPubKey asymmetric.IPubKey, pIncoming bool, pHash []byte) (bool, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	keyHash := getKeyHash(pPubKey, pHash)

	exist := false
	v, err := p.fDB.Get(keyHash)
	if err == nil {
		exist = true
		if v[0] == 1 {
			return false, ErrHashExist
		}
	}

	if !errors.Is(err, database.ErrNotFound) {
		return false, errors.Join(ErrGetHash, err)
	}

	state := byte(0)
	if pIncoming {
		state = byte(1)
	}

	if err := p.fDB.Set(keyHash, []byte{state}); err != nil {
		return false, errors.Join(ErrSetHash, err)
	}
	return exist, nil
}

func (p *sKeyValueDB) Load(pR IRelation, pStart, pEnd uint64) ([]IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	if pStart > pEnd {
		return nil, ErrStartGtEnd
	}

	size := p.getSize(pR)
	if pEnd > size {
		return nil, ErrEndGtSize
	}

	res := make([]IMessage, 0, pEnd-pStart)
	for i := pStart; i < pEnd; i++ {
		data, err := p.fDB.Get(getKeyMessageByEnum(pR, i))
		if err != nil {
			return nil, errors.Join(ErrGetMessage, err)
		}
		msg := LoadMessage(data)
		if msg == nil {
			return nil, ErrLoadMessage
		}
		res = append(res, msg)
	}

	return res, nil
}

func (p *sKeyValueDB) Push(pR IRelation, pMsg IMessage) error {
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
