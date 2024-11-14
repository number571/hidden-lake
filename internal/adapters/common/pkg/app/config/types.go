package config

import (
	"time"

	net_message "github.com/number571/go-peer/pkg/network/message"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() string
	GetConnection() IConnection
}

type IConfigSettings interface {
	net_message.ISettings

	GetWaitTime() time.Duration
}

type IConnection interface {
	GetHLTHost() string
	GetSrvHost() string
}
