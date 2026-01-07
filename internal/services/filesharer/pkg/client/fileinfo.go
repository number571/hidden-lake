package client

import "github.com/number571/go-peer/pkg/crypto/hashing"

var (
	_ IFileInfo = &sFileInfo{}
)

type sFileInfo struct {
	FName string `json:"name"`
	FHash string `json:"hash"`
	FSize uint64 `json:"size"`
}

func NewFileInfo(pName, pHash string, pSize uint64) IFileInfo {
	return &sFileInfo{
		FName: pName,
		FHash: pHash,
		FSize: pSize,
	}
}

func NewFileInfoFromBytes(pName string, b []byte) IFileInfo {
	return &sFileInfo{
		FName: pName,
		FHash: hashing.NewHasher(b).ToString(),
		FSize: uint64(len(b)),
	}
}

func (p *sFileInfo) GetName() string {
	return p.FName
}

func (p *sFileInfo) GetHash() string {
	return p.FHash
}

func (p *sFileInfo) GetSize() uint64 {
	return p.FSize
}
