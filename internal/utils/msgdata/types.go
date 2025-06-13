package msgdata

import (
	"html/template"
)

type IMessageBroker interface {
	Close(string) bool
	Produce(string, SMessage) bool
	Consume(string) (SMessage, bool)
}

type SSubscribe struct {
	FAddress string `json:"address"`
}

type SMessage struct {
	FTimestamp string        `json:"timestamp"`
	FTextData  template.HTML `json:"textdata"`
	FFileName  template.HTML `json:"filename"`
	FFileData  string        `json:"filedata"`
}
