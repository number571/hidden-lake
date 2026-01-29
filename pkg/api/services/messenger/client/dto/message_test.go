package dto

import (
	"testing"
	"time"
)

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
}
