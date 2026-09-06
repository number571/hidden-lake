package dto

import (
	"errors"
	"strings"

	"github.com/number571/go-peer/pkg/encoding"
)

func FileInfoListToString(list []IFileInfo) string {
	result := strings.Builder{}
	result.Grow(4096)

	result.WriteByte('[')

	if len(list) == 0 {
		result.WriteByte(']')
		return result.String()
	}

	for i := 0; i < len(list)-1; i++ {
		result.WriteString(list[i].ToString())
		result.WriteByte(',')
	}
	result.WriteString(list[len(list)-1].ToString())

	result.WriteByte(']')
	return result.String()
}

func LoadFileInfoList(pData interface{}) ([]IFileInfo, error) {
	var fileInfoListBytes []byte

	switch x := pData.(type) {
	case []byte:
		fileInfoListBytes = x
	case string:
		fileInfoListBytes = []byte(x)
	default:
		return nil, ErrUnknownType
	}

	list := []*sFileInfo{}
	if err := encoding.DeserializeJSON(fileInfoListBytes, &list); err != nil {
		return nil, errors.Join(ErrDecodeInfo, err)
	}

	for _, info := range list {
		if ok := isValidHexHash(info.FHash); !ok {
			return nil, ErrInvalidHash
		}
	}

	return fileInfoListToInterface(list), nil
}

func fileInfoListToInterface(list []*sFileInfo) []IFileInfo {
	result := make([]IFileInfo, 0, len(list))
	for _, v := range list {
		result = append(result, IFileInfo(v))
	}
	return result
}
