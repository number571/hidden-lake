package config

import (
	"github.com/number571/go-peer/pkg/client/message"
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
	message.ISettings
	net_message.ISettings

	GetMessagesCapacity() uint64
	GetRandMessageSizeBytes() uint64
	GetDatabaseEnabled() bool
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
	GetPPROF() string
}
