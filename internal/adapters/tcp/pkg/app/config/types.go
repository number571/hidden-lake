package config

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
	GetConnection() string
}

type IConfigSettings interface {
	net_message.ISettings
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
	GetPPROF() string
}
