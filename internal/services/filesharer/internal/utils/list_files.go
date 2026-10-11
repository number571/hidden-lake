package utils

import (
	"io/fs"
	"os"
)

func GetListFiles(pStgPath string, pPage uint64, pOffset uint64) ([]string, error) {
	stat, err := os.Stat(pStgPath)
	if os.IsNotExist(err) || !stat.IsDir() {
		return []string{}, nil
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

	result := make([]string, 0, pOffset)
	for i := (pPage * pOffset); i < uint64(len(files)); i++ {
		if fileReader == 0 {
			break
		}
		fileReader--

		fileName := files[i].Name()
		result = append(result, fileName)
	}

	return result, nil
}
