package database

import (
	"os"
	"testing"
	"time"

	message_dto "github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
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

	friend := "friend"

	timeNow := time.Now()
	msg, err := message_dto.NewMessageWithParams(message_dto.CIncomingState, tcBody, message_dto.NewMessageUID(), timeNow, false)
	if err != nil {
		t.Fatal(err)
	}

	_, err1 := db.Push(friend, msg)
	if err1 != nil {
		t.Fatal(err1)
	}

	size := db.Size(friend)
	if size != 1 {
		t.Fatal("size != 1")
	}

	loadMsg, err := db.LoadByEnum(friend, 0)
	if err != nil {
		t.Fatal(err)
	}

	if loadMsg.GetState() != message_dto.CIncomingState {
		t.Fatal("!msgs[0].IsIncoming()")
	}

	if loadMsg.GetBody() != tcBody {
		t.Fatal("!bytes.Equal(msgs[0].GetMessage(), []byte(tcBody))")
	}

	if loadMsg.GetTimestamp() != timeNow.Format(time.DateTime) {
		t.Fatal("msgs[0].GetTimestamp() != timeNow.Format(time.DateTime)")
	}
}
