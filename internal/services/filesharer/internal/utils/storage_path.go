package utils

import (
	"context"
	"path/filepath"

	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/pubkey"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
)

type iStgType uint

const (
	cStgPrivateType = iota
	cStgSharingType
)

func GetPrivateStoragePath(
	pCtx context.Context,
	pPathTo string,
	pHlkClient hlk_client.IClient,
	pAliasName string,
) (string, error) {
	return getStoragePath(pCtx, cStgPrivateType, pPathTo, pHlkClient, pAliasName, true)
}

func GetSharingStoragePath(
	pCtx context.Context,
	pPathTo string,
	pHlkClient hlk_client.IClient,
	pAliasName string,
	pPersonal bool,
) (string, error) {
	return getStoragePath(pCtx, cStgSharingType, pPathTo, pHlkClient, pAliasName, pPersonal)
}

func getStoragePath(
	pCtx context.Context,
	pStgType iStgType,
	pPathTo string,
	pHlkClient hlk_client.IClient,
	pAliasName string,
	pPersonal bool,
) (string, error) {
	if !pPersonal {
		return filepath.Join(pPathTo, hls_filesharer_settings.CPathSharingPublicSTG), nil
	}

	var directPath string
	switch pStgType {
	case cStgPrivateType:
		directPath = filepath.Join(pPathTo, hls_filesharer_settings.CPathPrivateSTG)
	case cStgSharingType:
		directPath = filepath.Join(pPathTo, hls_filesharer_settings.CPathSharingSTG)
	}

	fPubKey, err := pubkey.GetFriendPubKeyByAliasName(pCtx, pHlkClient, pAliasName)
	if err != nil {
		return "", err
	}

	return filepath.Join(directPath, fPubKey.GetHasher().ToString()), nil
}
