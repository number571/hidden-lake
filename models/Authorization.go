package models

type Authorization struct {
	Auth bool
	Login string
	Password []byte
}
