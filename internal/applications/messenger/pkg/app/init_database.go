package app

import (
	"fmt"
	"path/filepath"

	hlm_database "github.com/number571/hidden-lake/internal/applications/messenger/internal/database"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
)

func (p *sApp) initDatabase() error {
	db, err := hlm_database.NewKeyValueDB(filepath.Join(p.fPathTo, hlm_settings.CPathDB))
	if err != nil {
		return fmt.Errorf("open KV database: %w", err)
	}
	p.fDatabase = db
	return nil
}
