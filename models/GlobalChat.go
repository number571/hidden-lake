package models

type GlobalChat struct {
	Head GlobalChatHead `json:"head"`
	Body GlobalChatBody `json:"body"`
}

type GlobalChatHead struct {
	Founder string           `json:"founder"`
	Option  string           `json:"option"`
	Sender  GlobalChatSender `json:"sender"`
}

type GlobalChatBody struct {
	Data string         `json:"data"`
	Desc GlobalChatDesc `json:"desc"`
}

type GlobalChatSender struct {
	Public   string `json:"public_key"`
	Hashname string `json:"hashname"`
}

type GlobalChatDesc struct {
	Rand string `json:"rand"`
	Hash string `json:"hash"`
	Sign string `json:"sign"`
}
