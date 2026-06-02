package config

import (
	"time"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
	GetConnection() IConnection
	GetEndpoints() []string
}

type IConfigSettings interface {
	layer1.ISettings

	GetMessageSizeBytes() uint64
	GetDatabaseEnabled() bool

	GetWatchPeriod() time.Duration
	GetMaxDelayTime() time.Duration
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
}

type IConnection interface {
	GetDevPath() string
	GetChannel() uint8
}

type IAddress interface {
	GetExternal() string
	GetInternal() string
}
