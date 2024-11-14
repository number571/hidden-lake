package config

import (
	"time"

	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IConfig interface {
	GetSettings() IConfigSettings
	GetAddress() IAddress
	GetLogging() logger.ILogging
}

type IConfigSettings interface {
	GetExecTimeout() time.Duration
	GetPassword() string
}

type IAddress interface {
	GetIncoming() string
	GetPPROF() string
}
