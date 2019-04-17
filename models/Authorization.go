package models

type Authorization struct {
	Auth bool
	Hash string
	Login string
	Password []byte
}
