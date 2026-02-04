package utils

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

func GetFileInfoList(pStgPath string, pPage uint64, pOffset uint64) (dto.IFileInfoList, error) {
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

	return dto.LoadFileInfoList(result)
}
