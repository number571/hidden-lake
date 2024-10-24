package config

import (
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateConnections([]string) error
	UpdateFriends(map[string]asymmetric.IPubKey) error
}

type IConfigSettings interface {
	net_message.ISettings

	GetMessageSizeBytes() uint64
	GetFetchTimeoutMS() uint64
	GetQueuePeriodMS() uint64
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetFriends() map[string]asymmetric.IPubKey
	GetConnections() []string
	GetService(string) (IService, bool)
}

type IService interface {
	GetHost() string
}

type IAddress interface {
	GetTCP() string
	GetHTTP() string
	GetPPROF() string
}
