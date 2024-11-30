package utils

import "html/template"

type SSubscribe struct {
	FAddress string `json:"address"`
}

type SMessage struct {
	FTimestamp string        `json:"timestamp"`
	FTextData  template.HTML `json:"textdata"`
	FFileName  template.HTML `json:"filename"`
	FFileData  string        `json:"filedata"`
}
