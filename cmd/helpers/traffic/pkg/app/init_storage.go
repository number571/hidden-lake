package app

import (
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/cmd/helpers/traffic/internal/cache"
	"github.com/number571/hidden-lake/cmd/helpers/traffic/internal/storage"
)

func (p *sApp) initStorage(pDatabase database.IKVDatabase) {
	cfgSettings := p.fConfig.GetSettings()
	p.fStorage = storage.NewMessageStorage(
		cfgSettings,
		pDatabase,
		cache.NewLRUCache(cfgSettings.GetMessagesCapacity()),
	)
}
