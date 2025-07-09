package app

import (
	"fmt"
	"path/filepath"

	hls_messenger_database "github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
)

func (p *sApp) initDatabase() error {
	db, err := hls_messenger_database.NewKeyValueDB(filepath.Join(p.fPathTo, hls_messenger_settings.CPathDB))
	if err != nil {
		return fmt.Errorf("open KV database: %w", err)
	}
	p.fDatabase = db
	return nil
}
