package utils

import (
	"path/filepath"

	"github.com/number571/hidden-lake/internal/utils/chars"
)

func FileNameIsInvalid(pFilename string) bool {
	switch {
	case pFilename == "":
		return true
	case chars.HasNotGraphicCharacters(pFilename):
		return true
	case pFilename != filepath.Base(pFilename):
		return true
	}
	return false
}
