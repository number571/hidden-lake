package utils

import (
	"io/fs"
	"os"
	"path/filepath"

	fileinfo "github.com/number571/hidden-lake/pkg/api/services/filesharer/client/dto"
)

func GetListFileInfo(pStgPath string, pPage uint64, pOffset uint64) (fileinfo.IFileInfoList, error) {
	fileReader := pOffset

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

	result := make([]fileinfo.IFileInfo, 0, pOffset)
	for i := (pPage * pOffset); i < uint64(len(files)); i++ {
		if fileReader == 0 {
			break
		}
		fileReader--

		fileName := files[i].Name()
		fullPath := filepath.Join(pStgPath, fileName)

		info, err := fileinfo.NewFileInfo(fullPath)
		if err != nil {
			return nil, err
		}

		result = append(result, info)
	}

	return fileinfo.LoadFileInfoList(result)
}
