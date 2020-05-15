package models

type GroupChat struct {
	Head GroupChatHead `json:"head"`
	Body GroupChatBody `json:"body"`
}

type GroupChatHead struct {
	Founder string          `json:"founder"`
	Option  string          `json:"option"`
	Sender  GroupChatSender `json:"sender"`
}

type GroupChatBody struct {
	Data string        `json:"data"`
	Desc GroupChatDesc `json:"desc"`
}

type GroupChatSender struct {
	Public   string `json:"public_key"`
	Hashname string `json:"hashname"`
}

type GroupChatDesc struct {
	Rand string `json:"rand"`
	Hash string `json:"hash"`
	Sign string `json:"sign"`
}
