package sealed

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/number571/go-peer/pkg/storage/database"
	"go.etcd.io/bbolt"
)

const (
	cSeenKeySize     = 32
	cSeenBucket      = "_SEEN_"
	cSeenOpenTimeout = time.Second
	// cSeenCheckKey holds a key-derived tag so a wrong seen_key or a file
	// swapped in by an attacker is detected instead of silently passing every
	// anti-replay lookup. Short literal, never collides with a 32-byte HMAC key.
	cSeenCheckKey = "\x00check"
	cSeenCheckTag = "seen-db-check-v1"
)

var (
	_ database.IKVDatabase = &sSeenDB{}
)

// sSeenDB stores records under HMAC-SHA256(seen_key, key): the file on disk
// does not reveal which network messages (hashes) the node has seen.
// bbolt is opened directly (not via go-peer) to bound the wait for a file lock.
type sSeenDB struct {
	fKey []byte
	fDB  *bbolt.DB
}

// NewSeenDB opens (creates) the bbolt file pPath. The parent directory
// must exist; a file locked by another process is an error after ~1s.
func NewSeenDB(pPath string, pKey []byte) (database.IKVDatabase, error) {
	if len(pKey) != cSeenKeySize {
		return nil, ErrInvalidSeenKey
	}
	fi, err := os.Stat(filepath.Dir(pPath))
	if err != nil {
		return nil, errors.Join(ErrOpenSeenDB, ErrSeenDBDirNotFound, err)
	}
	if !fi.IsDir() {
		return nil, errors.Join(ErrOpenSeenDB, ErrSeenDBDirNotFound)
	}
	db, err := bbolt.Open(pPath, 0600, &bbolt.Options{Timeout: cSeenOpenTimeout})
	if err != nil {
		return nil, errors.Join(ErrOpenSeenDB, err)
	}
	p := &sSeenDB{
		fKey: append([]byte{}, pKey...),
		fDB:  db,
	}
	if err := p.verifyKey(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return p, nil
}

// verifyKey writes the key-check tag on a fresh file, or rejects a file whose
// tag was written under a different key. A rollback to an older genuine copy
// (same key) is not detected — see the accepted risks.
func (p *sSeenDB) verifyKey() error {
	want := p.hashKey([]byte(cSeenCheckTag))
	err := p.fDB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(cSeenBucket))
		if err != nil {
			return err
		}
		got := bucket.Get([]byte(cSeenCheckKey))
		if got == nil {
			return bucket.Put([]byte(cSeenCheckKey), want)
		}
		if !hmac.Equal(got, want) {
			return ErrSeenDBKeyMismatch
		}
		return nil
	})
	if err != nil {
		return errors.Join(ErrOpenSeenDB, err)
	}
	return nil
}

func (p *sSeenDB) Set(pKey []byte, pValue []byte) error {
	err := p.fDB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(cSeenBucket))
		if err != nil {
			return err
		}
		return bucket.Put(p.hashKey(pKey), pValue)
	})
	if err != nil {
		return errors.Join(database.ErrSetValueDB, err)
	}
	return nil
}

// Get returns database.ErrNotFound (as go-peer QB node expects) for a missing key.
func (p *sSeenDB) Get(pKey []byte) ([]byte, error) {
	var value []byte
	err := p.fDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(cSeenBucket))
		if bucket == nil {
			return database.ErrNotFound
		}
		val := bucket.Get(p.hashKey(pKey))
		if val == nil {
			return database.ErrNotFound
		}
		value = append([]byte{}, val...)
		return nil
	})
	if err != nil {
		return nil, errors.Join(database.ErrGetValueDB, err)
	}
	return value, nil
}

func (p *sSeenDB) Del(pKey []byte) error {
	err := p.fDB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(cSeenBucket))
		if bucket == nil {
			return nil
		}
		return bucket.Delete(p.hashKey(pKey))
	})
	if err != nil {
		return errors.Join(database.ErrDelValueDB, err)
	}
	return nil
}

func (p *sSeenDB) Close() error {
	if err := p.fDB.Close(); err != nil {
		return errors.Join(database.ErrCloseDB, err)
	}
	return nil
}

func (p *sSeenDB) hashKey(pKey []byte) []byte {
	h := hmac.New(sha256.New, p.fKey)
	_, _ = h.Write(pKey)
	return h.Sum(nil)
}
