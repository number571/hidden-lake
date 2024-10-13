package app

import (
	"fmt"
	"path/filepath"

	"github.com/number571/go-peer/pkg/storage/database"
	hlt_database "github.com/number571/hidden-lake/cmd/helpers/traffic/internal/database"
	hlt_settings "github.com/number571/hidden-lake/cmd/helpers/traffic/pkg/settings"
)

func (p *sApp) initDatabase() error {
	if !p.fConfig.GetSettings().GetDatabaseEnabled() {
		p.fDatabase = hlt_database.NewVoidKVDatabase()
		return nil
	}
	db, err := database.NewKVDatabase(filepath.Join(p.fPathTo, hlt_settings.CPathDB))
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}
	p.fDatabase = db
	return nil
}
