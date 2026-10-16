package database

import (
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/storage/database"
	message "github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
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

func (p *sKeyValueDB) Size(pR string) uint64 {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.getSize(pR)
}

func (p *sKeyValueDB) LoadByEnum(pR string, pNID uint64) (message.IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	data, err := p.fDB.Get(getKeyMessageByEnum(pR, pNID))
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errors.Join(ErrNotFoundMessage, err)
		}
		return nil, errors.Join(ErrGetMessage, err)
	}

	msg, err := message.LoadMessage(data)
	if err != nil {
		return nil, errors.Join(ErrLoadMessage, err)
	}

	return msg, nil
}

func (p *sKeyValueDB) LoadByUID(pR string, pUID string) (message.IMessage, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	return p.loadByUID(pR, pUID)
}

func (p *sKeyValueDB) Push(pR string, pMsg message.IMessage) (uint64, error) {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	msgUID := pMsg.GetHead().GetUID()

	if _, err := p.loadIndexByUID(pR, msgUID); err == nil {
		return 0, ErrExistMessage
	}

	size := p.getSize(pR)

	currPtr := encoding.Uint64ToBytes(size)
	nextPtr := encoding.Uint64ToBytes(size + 1)

	if err := p.fDB.Set(getKeyMessageByEnum(pR, size), pMsg.ToBytes()); err != nil {
		return 0, errors.Join(ErrSetMessage, err)
	}
	if err := p.fDB.Set(getKeyMessageByUID(pR, msgUID), currPtr[:]); err != nil {
		return 0, errors.Join(ErrSetMessage, err)
	}

	if err := p.fDB.Set(getKeySize(pR), nextPtr[:]); err != nil {
		return 0, errors.Join(ErrSetSizeMessage, err)
	}

	return size, nil
}

func (p *sKeyValueDB) Recv(pR string, pUID string) error {
	p.fMutex.Lock()
	defer p.fMutex.Unlock()

	index, err := p.loadIndexByUID(pR, pUID)
	if err != nil {
		return errors.Join(ErrGetIndexMessage, err)
	}

	key := getKeyMessageByEnum(pR, index)
	data, err := p.fDB.Get(key)
	if err != nil {
		return errors.Join(ErrGetMessage, err)
	}

	msg, err := message.LoadMessage(data)
	if err != nil {
		return errors.Join(ErrLoadMessage, err)
	}

	msgHead := msg.GetHead()

	switch msgHead.GetState() {
	case message.CReceivedState:
		// pass
	case message.CSendingState:
		newMsg := message.NewMessage(
			message.NewMessageHead(message.CReceivedState, msgHead.GetUID()),
			msg.GetBody(),
		)
		if err := p.fDB.Set(key, newMsg.ToBytes()); err != nil {
			return errors.Join(ErrSetMessage, err)
		}
	default:
		return ErrInvalidState
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

func (p *sKeyValueDB) loadByUID(pR string, pUID string) (message.IMessage, error) {
	index, err := p.loadIndexByUID(pR, pUID)
	if err != nil {
		return nil, errors.Join(ErrGetIndexMessage, err)
	}

	key := getKeyMessageByEnum(pR, index)
	data, err := p.fDB.Get(key)
	if err != nil {
		return nil, errors.Join(ErrGetMessage, err)
	}

	msg, err := message.LoadMessage(data)
	if err != nil {
		return nil, errors.Join(ErrLoadMessage, err)
	}

	return msg, nil
}

func (p *sKeyValueDB) loadIndexByUID(pR string, pUID string) (uint64, error) {
	key := getKeyMessageByUID(pR, pUID)

	ptr, err := p.fDB.Get(key)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return 0, errors.Join(ErrNotFoundMessage, err)
		}
		return 0, errors.Join(ErrGetMessage, err)
	}

	b := [encoding.CSizeUint64]byte{}
	copy(b[:], ptr)

	return encoding.BytesToUint64(b), nil
}

func (p *sKeyValueDB) getSize(pR string) uint64 {
	data, err := p.fDB.Get(getKeySize(pR))
	if err != nil {
		return 0
	}

	res := [encoding.CSizeUint64]byte{}
	copy(res[:], data)
	return encoding.BytesToUint64(res)
}
