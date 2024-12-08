package config

import (
	net_message "github.com/number571/go-peer/pkg/network/message"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateConnections([]string) error
}

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
	GetEndpoint() string
	GetConnections() []string
}

type IConfigSettings interface {
	net_message.ISettings

	GetMessageSizeBytes() uint64
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
	GetPPROF() string
}
