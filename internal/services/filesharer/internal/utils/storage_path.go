package utils

import (
	"path/filepath"

	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
)

type iStgType uint

const (
	cStgPrivateType = iota
	cStgSharingType
)

func GetPrivateStoragePath(
	pPathTo string,
	pAliasName string,
) string {
	return getStoragePath(cStgPrivateType, pPathTo, pAliasName, false)
}

func GetSharingStoragePath(
	pPathTo string,
	pAliasName string,
	pPublic bool,
) string {
	return getStoragePath(cStgSharingType, pPathTo, pAliasName, pPublic)
}

func getStoragePath(
	pStgType iStgType,
	pPathTo string,
	pAliasName string,
	pPublic bool,
) string {
	if pPublic {
		return filepath.Join(pPathTo, hls_filesharer_settings.CPathSharingPublicSTG)
	}

	var directPath string
	switch pStgType {
	case cStgPrivateType:
		directPath = filepath.Join(pPathTo, hls_filesharer_settings.CPathPrivateSTG)
	case cStgSharingType:
		directPath = filepath.Join(pPathTo, hls_filesharer_settings.CPathSharingSTG)
	}

	return filepath.Join(directPath, pAliasName)
}
