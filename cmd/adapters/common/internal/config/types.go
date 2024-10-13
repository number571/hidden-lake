package config

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
	logger "github.com/number571/hidden-lake/internal/logger/std"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() string
	GetConnection() IConnection
}

type IConfigSettings interface {
	net_message.ISettings

	GetWaitTimeMS() uint64
}

type IConnection interface {
	GetHLTHost() string
	GetSrvHost() string
}
