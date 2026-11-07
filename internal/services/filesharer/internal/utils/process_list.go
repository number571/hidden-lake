package utils

import (
	"strings"

	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

func DownloadProcessKeysToString(keys []dto.IDownloadProcessKey) string {
	result := strings.Builder{}
	result.Grow(4096)

	result.WriteByte('[')

	if len(keys) == 0 {
		result.WriteByte(']')
		return result.String()
	}

	for i := 0; i < len(keys)-1; i++ {
		result.WriteString(keys[i].ToString())
		result.WriteByte(',')
	}
	result.WriteString(keys[len(keys)-1].ToString())

	result.WriteByte(']')
	return result.String()
}
