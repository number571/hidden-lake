package config

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
	logger "github.com/number571/hidden-lake/internal/modules/logger/std"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
	GetProducers() []string
	GetConsumers() []string
}

type IConfigSettings interface {
	net_message.ISettings

	GetMessagesCapacity() uint64
}

type IAddress interface {
	GetHTTP() string
	GetPPROF() string
}
