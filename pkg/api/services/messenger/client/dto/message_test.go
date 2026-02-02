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

	isIncoming := true
	timeNow := time.Now()
	msgBody := "hello, world!"
	msg := NewMessage(isIncoming, msgBody, timeNow)

	loadMsg, err := LoadMessage(msg.ToBytes())
	if err != nil {
		t.Fatal(err)
	}
	if loadMsg.IsIncoming() != isIncoming {
		t.Fatal("is_incoming is invalid")
	}
	if loadMsg.GetMessage() != msgBody {
		t.Fatal("message is invalid")
	}
	if loadMsg.GetTimestamp() != timeNow.Format(time.DateTime) {
		t.Fatal("timestamp is invalid")
	}

	y, err := ParseTimestamp(timeNow.Format(time.DateTime))
	if err != nil {
		panic(err)
	}
	if x, err := ParseTimestamp(loadMsg.GetTimestamp()); err != nil || !x.Equal(y) {
		t.Fatal("parse message timestamp")
	}

	loadMsg2, err := LoadMessage(msg.ToString())
	if err != nil {
		t.Fatal(err)
	}
	if loadMsg2.IsIncoming() != isIncoming {
		t.Fatal("is_incoming is invalid")
	}
	if loadMsg2.GetMessage() != msgBody {
		t.Fatal("message is invalid")
	}
	if loadMsg2.GetTimestamp() != timeNow.Format(time.DateTime) {
		t.Fatal("timestamp is invalid")
	}
	if x, err := ParseTimestamp(loadMsg2.GetTimestamp()); err != nil || !x.Equal(y) {
		t.Fatal("parse message timestamp")
	}

	if _, err := LoadMessage(1); err == nil {
		t.Fatal("success load invalid type")
	}
	if _, err := LoadMessage([]byte{1}); err == nil {
		t.Fatal("success load invalid json")
	}
}
