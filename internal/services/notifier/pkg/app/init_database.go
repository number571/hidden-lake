package app

import (
	"fmt"
	"path/filepath"

	hls_notifier_database "github.com/number571/hidden-lake/internal/services/notifier/internal/database"
	hls_notifier_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
)

func (p *sApp) initDatabase() error {
	db, err := hls_notifier_database.NewKeyValueDB(filepath.Join(p.fPathTo, hls_notifier_settings.CPathDB))
	if err != nil {
		return fmt.Errorf("open KV database: %w", err)
	}
	p.fDatabase = db
	return nil
}
