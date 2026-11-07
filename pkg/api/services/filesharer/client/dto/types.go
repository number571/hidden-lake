package dto

import "github.com/number571/go-peer/pkg/types"

const (
	CTrimHashSize = 20 // hex encode
)

type IFileInfo interface {
	types.IConverter

	GetName() string
	GetHash() string
	GetSize() uint64
}

type IDownloadProcess interface {
	types.IConverter

	IDownloadProcessKey
	IDownloadProcessValue
}

type IDownloadProcessKey interface {
	types.IConverter

	GetFriend() string
	GetFileName() string
	GetIsPersonal() bool
}

type IDownloadProcessValue interface {
	types.IConverter

	GetIncIndex() uint64
	GetDownload() uint64
	GetFileSize() uint64
}
