package config

import (
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
	GetMessageSizeBytes() uint64
	GetDatabaseEnabled() bool
}

type IConnection interface {
	GetDevPath() string
	GetChannel() uint8
}

type IAddress interface {
	GetExternal() string
	GetInternal() string
}
