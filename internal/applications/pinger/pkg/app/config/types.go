package config

import (
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IConfig interface {
	GetSettings() IConfigSettings
	GetAddress() IAddress
	GetLogging() logger.ILogging
}

type IConfigSettings interface {
	GetResponseMessage() string
}

type IAddress interface {
	GetExternal() string
	GetPPROF() string
}
