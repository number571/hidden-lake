package dto

import (
	"errors"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	_ IFileInfoList = sFileInfoList{}
)

type sFileInfoList []*sFileInfo

func LoadFileInfoList(pData interface{}) (IFileInfoList, error) {
	var fileInfoListBytes []byte

	switch x := pData.(type) {
	case []byte:
		fileInfoListBytes = x
	case string:
		fileInfoListBytes = []byte(x)
	case IFileInfo:
		info := &sFileInfo{FName: x.GetName(), FHash: x.GetHash(), FSize: x.GetSize()}
		return sFileInfoList{info}, nil
	case []IFileInfo:
		list := make(sFileInfoList, 0, len(x))
		for _, f := range x {
			info := &sFileInfo{FName: f.GetName(), FHash: f.GetHash(), FSize: f.GetSize()}
			list = append(list, info)
		}
		return list, nil
	default:
		return nil, ErrUnknownType
	}

	list := sFileInfoList{}
	if err := encoding.DeserializeJSON(fileInfoListBytes, &list); err != nil {
		return nil, errors.Join(ErrDecodeInfo, err)
	}

	for _, info := range list {
		if ok := isValidHexHash(info.FHash); !ok {
			return nil, ErrInvalidHash
		}
	}

	return list, nil
}

func (p sFileInfoList) ToString() string {
	return string(p.ToBytes())
}

func (p sFileInfoList) ToBytes() []byte {
	return encoding.SerializeJSON(p)
}

func (p sFileInfoList) GetList() []IFileInfo {
	result := make([]IFileInfo, 0, len(p))
	for _, info := range p {
		result = append(result, info)
	}
	return result
}
