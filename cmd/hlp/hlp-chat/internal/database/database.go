package database

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/symmetric"
	gp_database "github.com/number571/go-peer/pkg/storage/database"
)

type sDatabase struct {
	fMtx *sync.Mutex
	fKVD gp_database.IKVDatabase
	fKey [3][]byte
}

func NewDatabase(pPath string, pKey [3][]byte) (IDatabase, error) {
	for _, k := range pKey {
		if len(k) != symmetric.CCipherKeySize {
			return nil, ErrKeySize
		}
	}
	kvDB, err := gp_database.NewKVDatabase(pPath)
	if err != nil {
		return nil, err
	}
	return &sDatabase{
		fMtx: &sync.Mutex{},
		fKVD: kvDB,
		fKey: pKey,
	}, nil
}

func (p *sDatabase) GetOrigin() gp_database.IKVDatabase {
	return p.fKVD
}

func (p *sDatabase) Insert(pChannel asymmetric.IPubKey, pMsg SMessage) error {
	p.fMtx.Lock()
	defer p.fMtx.Unlock()

	count, err := p.getCount(pChannel)
	if err != nil {
		return err
	}

	if err := p.fKVD.Set(p.keyGetMsg(pChannel, count), p.messageToBytes(pMsg)); err != nil {
		return errors.Join(ErrSetCount, err)
	}
	if err := p.fKVD.Set(p.keyCountMsgs(pChannel), []byte(fmt.Sprintf("%d", count+1))); err != nil { // nolint: perfsprint
		return errors.Join(ErrSetMessage, err)
	}

	return nil
}

func (p *sDatabase) Select(pChannel asymmetric.IPubKey, pN uint64) ([]SMessage, error) {
	p.fMtx.Lock()
	defer p.fMtx.Unlock()

	count, err := p.getCount(pChannel)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	readUintil := uint64(0)
	if count > pN {
		readUintil = count - pN
	}

	msgs := make([]SMessage, 0, pN)
	for i := int64(count - 1); i >= int64(readUintil); i-- { // nolint: gosec
		msgBytes, err := p.fKVD.Get(p.keyGetMsg(pChannel, uint64(i))) // nolint: gosec
		if err != nil {
			return nil, errors.Join(ErrGetMessage, err)
		}
		msg, err := p.bytesToMessage(msgBytes)
		if err != nil {
			return nil, errors.Join(ErrDecodeMsg, err)
		}
		msgs = append(msgs, msg)
	}

	slices.Reverse(msgs)
	return msgs, nil
}

func (p *sDatabase) getCount(pChannel asymmetric.IPubKey) (uint64, error) {
	countBytes, err := p.fKVD.Get(p.keyCountMsgs(pChannel))
	if err != nil {
		if !errors.Is(err, gp_database.ErrNotFound) {
			return 0, errors.Join(ErrGetCount, err)
		}
		countBytes = []byte("0")
		if err := p.fKVD.Set(p.keyCountMsgs(pChannel), countBytes); err != nil {
			return 0, errors.Join(ErrSetCount, err)
		}
	}
	count, err := strconv.ParseUint(string(countBytes), 10, 64)
	if err != nil {
		return 0, errors.Join(ErrParseCount, err)
	}
	return count, nil
}

func (p *sDatabase) messageToBytes(pMsg SMessage) []byte {
	msgBytes := make([]byte, 0, 256)
	msgBytes = append(msgBytes, []byte(pMsg.FSendTime.Format(time.DateTime))...)
	msgBytes = append(msgBytes, []byte(pMsg.FSender)...)
	msgBytes = append(msgBytes, []byte(pMsg.FMessage)...)
	return p.encryptBytes(msgBytes)
}

func (p *sDatabase) bytesToMessage(pEncBytes []byte) (SMessage, error) {
	msgBytes, err := p.decryptBytes(pEncBytes)
	if err != nil {
		return SMessage{}, err
	}
	dateTimeSize := len(time.DateTime)
	if len(msgBytes) < ed25519.PublicKeySize+dateTimeSize {
		return SMessage{}, ErrMsgSize
	}
	sendTime, err := time.Parse(time.DateTime, string(msgBytes[:dateTimeSize]))
	if err != nil {
		return SMessage{}, err
	}
	return SMessage{
		FSendTime: sendTime,
		FSender:   ed25519.PublicKey(msgBytes[dateTimeSize : dateTimeSize+ed25519.PublicKeySize]),
		FMessage:  string(msgBytes[dateTimeSize+ed25519.PublicKeySize:]),
	}, nil
}
