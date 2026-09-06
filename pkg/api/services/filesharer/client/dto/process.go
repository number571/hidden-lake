package dto

import (
	"errors"

	"github.com/number571/go-peer/pkg/encoding"
)

type sDownloadProcess struct {
	*sDownloadProcessKey   `json:"key"`
	*sDownloadProcessValue `json:"value"`
}

type sDownloadProcessKey struct {
	FFriend   string `json:"friend"`
	FName     string `json:"name"`
	FPersonal bool   `json:"personal"`
}

type sDownloadProcessValue struct {
	FIncIndex uint64 `json:"incindex"`
	FDownload uint64 `json:"download"`
	FFileSize uint64 `json:"filesize"`
}

func LoadDownloadProcess(pData interface{}) (IDownloadProcess, error) {
	var downloadProcessBytes []byte

	switch x := pData.(type) {
	case []byte:
		downloadProcessBytes = x
	case string:
		downloadProcessBytes = []byte(x)
	default:
		return nil, ErrUnknownType
	}

	downloadProcess := &sDownloadProcess{}
	if err := encoding.DeserializeJSON(downloadProcessBytes, downloadProcess); err != nil {
		return nil, errors.Join(ErrDecodeInfo, err)
	}

	return downloadProcess, nil
}

func NewDownloadProcess(k IDownloadProcessKey, v IDownloadProcessValue) IDownloadProcess {
	return &sDownloadProcess{
		k.(*sDownloadProcessKey),
		v.(*sDownloadProcessValue),
	}
}

func NewDownloadProcessValue(incIndex, download, fileSize uint64) IDownloadProcessValue {
	return &sDownloadProcessValue{
		FIncIndex: incIndex,
		FDownload: download,
		FFileSize: fileSize,
	}
}

func NewDownloadProcessKey(friend, filename string, isPersonal bool) IDownloadProcessKey {
	return &sDownloadProcessKey{
		FFriend:   friend,
		FName:     filename,
		FPersonal: isPersonal,
	}
}

func (p *sDownloadProcess) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p *sDownloadProcess) ToString() string {
	return string(p.ToBytes())
}

func (p *sDownloadProcessValue) GetIncIndex() uint64 {
	return p.FIncIndex
}

func (p *sDownloadProcessValue) GetDownload() uint64 {
	return p.FDownload
}

func (p *sDownloadProcessValue) GetFileSize() uint64 {
	return p.FFileSize
}

func (p *sDownloadProcessKey) GetFriend() string {
	return p.FFriend
}

func (p *sDownloadProcessKey) GetFileName() string {
	return p.FName
}

func (p *sDownloadProcessKey) GetIsPersonal() bool {
	return p.FPersonal
}
