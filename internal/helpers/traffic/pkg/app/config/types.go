package config

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IConfig interface {
	GetSettings() IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetConnections() []string
	GetConsumers() []string
}

type IConfigSettings interface {
	net_message.ISettings

	GetMessageSizeBytes() uint64
	GetMessagesCapacity() uint64
	GetDatabaseEnabled() bool
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
	GetPPROF() string
}
