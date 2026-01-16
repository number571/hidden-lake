package config

import (
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IConfig interface {
	GetSettings() IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetConnection() string
}

type IConfigSettings interface {
	GetRetryNum() uint64
	GetPageOffset() uint64
}

type IAddress interface {
	GetInternal() string
	GetExternal() string
}
