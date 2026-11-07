package app

import (
	"fmt"
	"path/filepath"

	hls_publisher_database "github.com/number571/hidden-lake/internal/services/publisher/internal/database"
	hls_publisher_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
)

func (p *sApp) initDatabase() error {
	db, err := hls_publisher_database.NewKeyValueDB(filepath.Join(p.fPathTo, hls_publisher_settings.CPathDB))
	if err != nil {
		return fmt.Errorf("open KV database: %w", err)
	}
	p.fDatabase = db
	return nil
}
