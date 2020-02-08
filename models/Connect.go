package models

type Connect struct {
	Hidden     bool   `json:"hidden"`
	Connected  bool   `json:"connected"`
	Address    string `json:"address"`
	Hashname   string `json:"hashname"`
	ThrowNode  string `json:"thrownode"`
	PublicKey  string `json:"public_key"`
}
