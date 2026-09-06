package utils

import (
	"strings"

	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

func DownloadProcessListToString(list []dto.IDownloadProcess) string {
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
