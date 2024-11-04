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
	UpdateFriends(map[string]asymmetric.IPubKey) error
}

type IConfig interface {
	GetLogging() logger.ILogging
	GetSettings() IConfigSettings

	GetAddress() IAddress
	GetFriends() map[string]asymmetric.IPubKey
}

type IConfigSettings interface {
	net_message.ISettings

	GetMessageSizeBytes() uint64
}

type IAddress interface {
	GetHTTP() string
	GetPPROF() string
}
