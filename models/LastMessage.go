package models

type LastMessage struct {
	Companion string  `json:"companion"`
	Message   Message `json:"message"`
}
