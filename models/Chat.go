package models

type Chat struct {
	Companion string    `json:"companion"`
	Messages  []Message `json:"messages"`
}
