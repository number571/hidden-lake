package app

import (
	"github.com/number571/go-peer/pkg/storage/database"
	"github.com/number571/hidden-lake/internal/adapters/tcp/internal/cache"
	"github.com/number571/hidden-lake/internal/adapters/tcp/internal/storage"
)

func (p *sApp) initStorage(pDatabase database.IKVDatabase) {
	cfgSettings := p.fWrapper.GetConfig().GetSettings()
	p.fStorage = storage.NewMessageStorage(
		cfgSettings,
		pDatabase,
		cache.NewLRUCache(cfgSettings.GetMessagesCapacity()),
	)
}
