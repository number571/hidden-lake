package dto

import (
	"errors"

	"github.com/number571/go-peer/pkg/encoding"
)

func LoadDownloadProcessList(pData interface{}) ([]IDownloadProcess, error) {
	var downloadProcessBytes []byte

	switch x := pData.(type) {
	case []byte:
		downloadProcessBytes = x
	case string:
		downloadProcessBytes = []byte(x)
	default:
		return nil, ErrUnknownType
	}

	list := []*sDownloadProcess{}
	if err := encoding.DeserializeJSON(downloadProcessBytes, &list); err != nil {
		return nil, errors.Join(ErrDecodeInfo, err)
	}

	return downloadProcessListToInterface(list), nil
}

func downloadProcessListToInterface(list []*sDownloadProcess) []IDownloadProcess {
	result := make([]IDownloadProcess, 0, len(list))
	for _, v := range list {
		result = append(result, IDownloadProcess(v))
	}
	return result
}
