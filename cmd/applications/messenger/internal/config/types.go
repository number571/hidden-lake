package config

import (
	"github.com/number571/hidden-lake/internal/language"
	logger "github.com/number571/hidden-lake/internal/logger/std"
)

type IWrapper interface {
	GetConfig() IConfig
	GetEditor() IEditor
}

type IEditor interface {
	UpdateLanguage(language.ILanguage) error
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetAddress() IAddress
	GetLogging() logger.ILogging
	GetConnection() string
}

type IConfigSettings interface {
	GetMessagesCapacity() uint64
	GetLanguage() language.ILanguage
}

type IAddress interface {
	GetInterface() string
	GetIncoming() string
	GetPPROF() string
}
