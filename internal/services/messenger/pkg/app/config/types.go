package config

import (
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IConfig interface {
	GetSettings() IConfigSettings
	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() string
}

type IConfigSettings interface {
	GetMessagesCapacity() uint64
}

type IAddress interface {
	GetInternal() string
	GetExternal() string
}
