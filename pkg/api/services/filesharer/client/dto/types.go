package dto

type IFileInfo interface {
	GetName() string
	GetHash() string
	GetSize() uint64
}
