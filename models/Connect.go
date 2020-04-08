package models

type Connect struct {
	Hidden       bool   `json:"hidden"`
	Connected    bool   `json:"connected"`
	InChat       bool   `json:"in_chat"`
	Address      string `json:"address"`
	Hashname     string `json:"hashname"`
	Public       string `json:"public_key"`
	ThrowClient  string `json:"throwclient"`
	Certificate  string `json:"certificate"`
}
