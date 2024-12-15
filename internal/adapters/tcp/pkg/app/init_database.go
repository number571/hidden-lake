package app

import (
	"fmt"
	"path/filepath"

	"github.com/number571/go-peer/pkg/storage/database"
	hla_tcp_database "github.com/number571/hidden-lake/internal/adapters/tcp/internal/database"
	hla_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"
)

func (p *sApp) initDatabase() error {
	if !p.fWrapper.GetConfig().GetSettings().GetDatabaseEnabled() {
		p.fDatabase = hla_tcp_database.NewVoidKVDatabase()
		return nil
	}
	db, err := database.NewKVDatabase(filepath.Join(p.fPathTo, hla_settings.CPathDB))
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}
	p.fDatabase = db
	return nil
}
