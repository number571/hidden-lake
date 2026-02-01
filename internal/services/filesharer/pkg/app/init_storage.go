package app

import (
	"os"
	"path/filepath"

	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
)

func (p *sApp) initStorage() error {
	stgPrivatePath := filepath.Join(p.fPathTo, hls_filesharer_settings.CPathPrivateSTG)
	if err := os.MkdirAll(stgPrivatePath, 0o700); err != nil { //nolint:gosec
		return err
	}
	stgSharingPath := filepath.Join(p.fPathTo, hls_filesharer_settings.CPathSharingSTG)
	if err := os.MkdirAll(stgSharingPath, 0o700); err != nil { //nolint:gosec
		return err
	}
	stgSharingPublicPath := filepath.Join(p.fPathTo, hls_filesharer_settings.CPathSharingPublicSTG)
	if err := os.MkdirAll(stgSharingPublicPath, 0o700); err != nil { //nolint:gosec
		return err
	}
	return nil
}
