package dto

import (
	"testing"
	"time"
)

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestMessage(t *testing.T) {
	t.Parallel()

	timeNow := time.Now()
	msgBody := "hello, world!"
	msg := NewMessage(
		NewMessageHead(CIncomingState, NewMessageUID()),
		NewMessageBody(msgBody, timeNow),
	)
	loadMsg, err := LoadMessage(msg.ToBytes())
	if err != nil {
		t.Fatal(err)
	}
	if loadMsg.GetHead().GetState() != CIncomingState {
		t.Fatal("is_incoming is invalid")
	}
	if loadMsg.GetBody().GetData() != msgBody {
		t.Fatal("message is invalid")
	}
	if loadMsg.GetBody().GetTimestamp() != timeNow.Format(time.DateTime) {
		t.Fatal("timestamp is invalid")
	}

	y, err := ParseTimestamp(timeNow.Format(time.DateTime))
	if err != nil {
		panic(err)
	}
	if x, err := ParseTimestamp(loadMsg.GetBody().GetTimestamp()); err != nil || !x.Equal(y) {
		t.Fatal("parse message timestamp")
	}

	loadMsg2, err := LoadMessage(msg.ToString())
	if err != nil {
		t.Fatal(err)
	}
	if loadMsg2.GetHead().GetState() != CIncomingState {
		t.Fatal("is_incoming is invalid")
	}
	if loadMsg2.GetBody().GetData() != msgBody {
		t.Fatal("message is invalid")
	}
	if loadMsg2.GetBody().GetTimestamp() != timeNow.Format(time.DateTime) {
		t.Fatal("timestamp is invalid")
	}
	if x, err := ParseTimestamp(loadMsg2.GetBody().GetTimestamp()); err != nil || !x.Equal(y) {
		t.Fatal("parse message timestamp")
	}

	if _, err := LoadMessage(1); err == nil {
		t.Fatal("success load invalid type")
	}
	if _, err := LoadMessage([]byte{1}); err == nil {
		t.Fatal("success load invalid json")
	}
}
