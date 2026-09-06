package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

func FileInfoListToString(list []dto.IFileInfo) string {
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

func GetFileInfoList(pStgPath string, pPage uint64, pOffset uint64) ([]dto.IFileInfo, error) {
	stat, err := os.Stat(pStgPath)
	if os.IsNotExist(err) || !stat.IsDir() {
		list, err := dto.LoadFileInfoList("[]")
		if err != nil {
			panic(err)
		}
		return list, nil
	}

	entries, err := os.ReadDir(pStgPath)
	if err != nil {
		return nil, err
	}

	files := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, e)
	}

	fileReader := pOffset

	result := make([]dto.IFileInfo, 0, pOffset)
	for i := (pPage * pOffset); i < uint64(len(files)); i++ {
		if fileReader == 0 {
			break
		}
		fileReader--

		fileName := files[i].Name()
		fullPath := filepath.Join(pStgPath, fileName)

		info, err := dto.NewFileInfo(fullPath)
		if err != nil {
			return nil, err
		}

		result = append(result, info)
	}

	return result, nil
}
