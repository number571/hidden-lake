package config

import (
	"github.com/number571/hidden-lake/internal/utils/language"
	logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateLanguage(language.ILanguage) error
	UpdateChannels([]string) error
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() string
	GetChannels() []string
}

type IConfigSettings interface {
	GetMessagesCapacity() uint64
	GetWorkSizeBits() uint64
	GetPowParallel() uint64
	GetLanguage() language.ILanguage
}

type IAddress interface {
	GetInternal() string
	GetExternal() string
}
