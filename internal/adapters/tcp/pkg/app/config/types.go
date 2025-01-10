package config

import (
	net_message "github.com/number571/go-peer/pkg/message/layer1"
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
	GetEndpoints() []string
	GetConnections() []string
}

type IConfigSettings interface {
	net_message.ISettings

	GetMessageSizeBytes() uint64
	GetDatabaseEnabled() bool
}

type IAddress interface {
	GetExternal() string
	GetInternal() string
	GetPPROF() string
}
