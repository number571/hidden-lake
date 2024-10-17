package database

import (
	"bytes"
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
)

const (
	tcBody = "hello, world!"
	tcPath = "database.db"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SDatabaseError{str}
	if err.Error() != errPrefix+str {
		t.Error("incorrect err.Error()")
		return
	}
}

func TestDatabase(t *testing.T) {
	t.Parallel()

	os.RemoveAll(tcPath)
	defer os.RemoveAll(tcPath)

	db, err := NewKeyValueDB(tcPath)
	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()

	iam := asymmetric.NewPrivKeyChain(
		asymmetric.NewKEncPrivKey(),
		asymmetric.NewSignPrivKey(),
	).GetPubKeyChain()
	friend := asymmetric.NewPrivKeyChain(
		asymmetric.NewKEncPrivKey(),
		asymmetric.NewSignPrivKey(),
	).GetPubKeyChain()

	rel := NewRelation(iam, friend)
	err1 := db.Push(rel, NewMessage(true, []byte(tcBody)))
	if err1 != nil {
		t.Error(err1)
		return
	}

	size := db.Size(rel)
	if size != 1 {
		t.Error("size != 1")
		return
	}

	msgs, err := db.Load(rel, 0, size)
	if err != nil {
		t.Error(err)
		return
	}

	if len(msgs) != 1 {
		t.Error("len(msgs) != 1")
		return
	}

	if !msgs[0].IsIncoming() {
		t.Error("!msgs[0].IsIncoming()")
		return
	}

	if !bytes.Equal(msgs[0].GetMessage(), []byte(tcBody)) {
		t.Error("!bytes.Equal(msgs[0].GetMessage(), []byte(tcBody))")
		return
	}
}
