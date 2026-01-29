package dto

import "github.com/number571/go-peer/pkg/types"

type IFileInfoList interface {
	types.IConverter

	GetList() []IFileInfo
}

type IFileInfo interface {
	types.IConverter

	GetName() string
	GetHash() string
	GetSize() uint64
}
