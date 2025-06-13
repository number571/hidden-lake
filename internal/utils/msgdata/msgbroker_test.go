package msgdata

import (
	"html/template"
	"testing"
	"time"
)

func TestMessageBroker(t *testing.T) {
	t.Parallel()

	addr := "address"
	msgData := template.HTML("msg_data")

	msgReceiver := NewMessageBroker()

	go func() {
		time.Sleep(100 * time.Millisecond)
		msgReceiver.Produce(addr, SMessage{FTextData: msgData})
		msgReceiver.Produce(addr, SMessage{FTextData: msgData})
		msgReceiver.Produce(addr, SMessage{FTextData: msgData})
	}()

	msg, ok := msgReceiver.Consume(addr)
	if !ok {
		t.Error("got not ok recv")
		return
	}
	if msg.FTextData != msgData {
		t.Error("msg.FTextData != msgData")
		return
	}

	msgReceiver.Close(addr)

	msgReceiver2 := NewMessageBroker()

	go func() {
		time.Sleep(500 * time.Millisecond)
		msgReceiver2.Produce(addr, SMessage{FTextData: msgData})
	}()

	go func() {
		if _, ok := msgReceiver2.Consume(addr); ok {
			t.Error("success consume with canceled")
			return
		}
	}()

	time.Sleep(100 * time.Millisecond)
	if _, ok := msgReceiver2.Consume(addr); !ok {
		t.Error("got not ok recv")
		return
	}
}

func TestReplaceTextToEmoji(t *testing.T) {
	t.Parallel()

	s := "hello :) !"
	got := ReplaceTextToEmoji(s)
	want := "hello 🙂 !"

	if got != want {
		t.Error("got incorrect replace text to emoji")
		return
	}
}

func TestReplaceTextToURLs(t *testing.T) {
	t.Parallel()

	s := "hello https://github.com/number571/github !"
	got := ReplaceTextToURLs(s)
	want := "hello <a style='background-color:#b9cdcf;color:black;' target='_blank' href='https://github.com/number571/github'>https://github.com/number571/github</a> !"

	if got != want {
		t.Error("got incorrect replace text to urls")
		return
	}
}
