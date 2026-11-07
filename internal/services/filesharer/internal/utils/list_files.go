package utils

import (
	"io/fs"
	"math"
	"os"
	"strings"
)

func GetListFiles(pStgPath string, pPage uint64, pOffset uint64) (uint64, []string, error) {
	stat, err := os.Stat(pStgPath)
	if os.IsNotExist(err) || !stat.IsDir() {
		return 0, []string{}, nil
	}

	entries, err := os.ReadDir(pStgPath)
	if err != nil {
		return 0, nil, err
	}

	files := make([]fs.DirEntry, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasPrefix(e.Name(), "~") {
			continue // pass file_info files if exists
		}
		files = append(files, e)
	}

	fileReader := pOffset
	filesCount := uint64(len(files))

	result := make([]string, 0, pOffset)
	for i := (pPage * pOffset); i < filesCount; i++ {
		if fileReader == 0 {
			break
		}
		fileReader--

		fileName := files[i].Name()
		result = append(result, fileName)
	}

	pagesCount := uint64(math.Ceil(float64(filesCount) / float64(pOffset)))
	return pagesCount, result, nil
}
