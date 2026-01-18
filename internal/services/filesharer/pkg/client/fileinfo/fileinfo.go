package fileinfo

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IFileInfo = &sFileInfo{}
)

type sFileInfo struct {
	FName string `json:"name"`
	FHash string `json:"hash"`
	FSize uint64 `json:"size"`
}

func LoadFileInfo(b []byte) (IFileInfo, error) {
	info := &sFileInfo{}
	if err := encoding.DeserializeJSON(b, info); err != nil {
		return nil, errors.Join(ErrDecodeInfo, err)
	}
	if ok := isValidHexHash(info.FHash); !ok {
		return nil, ErrInvalidHash
	}
	return info, nil
}

func LoadFileInfoList(b []byte) ([]IFileInfo, error) {
	list := []*sFileInfo{}
	if err := encoding.DeserializeJSON(b, &list); err != nil {
		return nil, errors.Join(ErrDecodeInfo, err)
	}
	fileInfos := make([]IFileInfo, 0, len(list))
	for _, info := range list {
		if ok := isValidHexHash(info.FHash); !ok {
			return nil, ErrInvalidHash
		}
		fileInfos = append(fileInfos, info)
	}
	return fileInfos, nil
}

func NewFileInfo(pName string) (IFileInfo, error) {
	hash, err := getFileHash(pName)
	if err != nil {
		return nil, err
	}
	size, err := getFileSize(pName)
	if err != nil {
		return nil, err
	}
	return &sFileInfo{
		FName: filepath.Base(pName),
		FHash: hash,
		FSize: size,
	}, nil
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

func isValidHexHash(hash string) bool {
	v, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}
	return len(v) == sha512.Size384
}

func getFileSize(filename string) (uint64, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return uint64(stat.Size()), nil //nolint:gosec
}

func getFileHash(filename string) (string, error) {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	h := sha512.New384()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return encoding.HexEncode(h.Sum(nil)), nil
}
