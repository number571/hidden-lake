package config

import (
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/message/layer1"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateFriends(map[string]asymmetric.IPubKey) error
}

type IConfigSettings interface {
	layer1.ISettings

	GetMessageSizeBytes() uint64
	GetFetchTimeout() time.Duration
	GetQueuePeriod() time.Duration

	GetPowParallel() uint64
	GetQBPConsumers() uint64
	GetQueueMainCap() uint64
	GetQueueRandCap() uint64
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetFriends() map[string]asymmetric.IPubKey
	GetEndpoints() []string
	GetService(string) (string, bool)
}

type IAddress interface {
	GetExternal() string
	GetInternal() string
}
