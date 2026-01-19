package database

import (
	"os"
	"testing"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/client/message"
)

const (
	tcBody = "hello, world!"
	tcPath = "database.db"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestDatabase(t *testing.T) {
	t.Parallel()

	_ = os.RemoveAll(tcPath)
	defer func() { _ = os.RemoveAll(tcPath) }()

	db, err := NewKeyValueDB(tcPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	iam := asymmetric.NewPrivKey().GetPubKey()
	friend := asymmetric.NewPrivKey().GetPubKey()

	rel := NewRelation(iam, friend)
	err1 := db.Push(rel, message.NewMessage(true, tcBody))
	if err1 != nil {
		t.Fatal(err1)
	}

	size := db.Size(rel)
	if size != 1 {
		t.Fatal("size != 1")
	}

	msgs, err := db.Load(rel, 0, size)
	if err != nil {
		t.Fatal(err)
	}

	if len(msgs) != 1 {
		t.Fatal("len(msgs) != 1")
	}

	if !msgs[0].IsIncoming() {
		t.Fatal("!msgs[0].IsIncoming()")
	}

	if msgs[0].GetMessage() != tcBody {
		t.Fatal("!bytes.Equal(msgs[0].GetMessage(), []byte(tcBody))")
	}
}
