package models

type EmailSaveOption bool
const (
	IsTempEmail EmailSaveOption = true
	IsPermEmail EmailSaveOption = false
)

type Email struct {
	Info  EmailInfo `json:"info"`
	Email EmailType `json:"email"`
}

type EmailInfo struct {
	Incoming  bool   `json:"incoming"`
	Time      string `json:"time"`
}

type EmailType struct {
	Head EmailHead `json:"head"`
	Body EmailBody `json:"body"`
}

type EmailHead struct {
	Sender     EmailSender `json:"sender"`
	Receiver   string      `json:"receiver"`   // hash(public receiver)
	Session    string      `json:"session"`    // encrypt[public receiver](session)
}

type EmailSender struct {
	Public   string `json:"public_key"` // public sender
	Hashname string `json:"hashname"`   // hash(public sender)
}

type EmailBody struct {
	Data string    `json:"data"` // encrypt[session](data)
	Desc EmailDesc `json:"desc"`
}

type EmailDesc struct {
	Rand       string `json:"rand"` // random(8)
	Hash       string `json:"hash"` // hash(data + random(8))
	Sign       string `json:"sign"` // sign[private sender](hash(data + random(8)))
	Nonce      uint64 `json:"nonce"`
	Difficulty uint8  `json:"difficulty"`
}
