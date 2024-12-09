package config

import (
	"time"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	net_message "github.com/number571/go-peer/pkg/network/message"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateFriends(map[string]asymmetric.IPubKey) error
}

type IConfigSettings interface {
	net_message.ISettings

	GetMessageSizeBytes() uint64
	GetFetchTimeout() time.Duration
	GetQueuePeriod() time.Duration
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetFriends() map[string]asymmetric.IPubKey
	GetAdapters() []string
	GetService(string) (string, bool)
}

type IAddress interface {
	GetExternal() string
	GetInternal() string
	GetPPROF() string
}
