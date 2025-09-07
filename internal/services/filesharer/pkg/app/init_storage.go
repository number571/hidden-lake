package app

import (
	"os"
	"path/filepath"

	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
)

func (p *sApp) initStorage() error {
	stgPath := filepath.Join(p.fPathTo, hls_filesharer_settings.CPathSTG)
	return os.MkdirAll(stgPath, 0o777) //nolint:gosec
}
