package utils

import "github.com/number571/go-peer/pkg/crypto/hashing"

type IFileInfo interface {
	GetName() string
	GetHash() string
	GetSize() uint64
}

var (
	_ IFileInfo = &SFileInfo{}
)

type SFileInfo struct {
	FName string `json:"name"`
	FHash string `json:"hash"`
	FSize uint64 `json:"size"`
}

func NewFileInfo(pName, pHash string, pSize uint64) IFileInfo {
	return &SFileInfo{
		FName: pName,
		FHash: pHash,
		FSize: pSize,
	}
}

func NewFileInfoFromBytes(pName string, b []byte) IFileInfo {
	return &SFileInfo{
		FName: pName,
		FHash: hashing.NewHasher(b).ToString(),
		FSize: uint64(len(b)),
	}
}

func (p *SFileInfo) GetName() string {
	return p.FName
}

func (p *SFileInfo) GetHash() string {
	return p.FHash
}

func (p *SFileInfo) GetSize() uint64 {
	return p.FSize
}
