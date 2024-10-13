package config

import (
	"github.com/number571/go-peer/pkg/client/message"
	net_message "github.com/number571/go-peer/pkg/network/message"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
}

type IConfigSettings interface {
	message.ISettings
	net_message.ISettings

	GetRandMessageSizeBytes() uint64
}

type IAddress interface {
	GetHTTP() string
	GetPPROF() string
}
