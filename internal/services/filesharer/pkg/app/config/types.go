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
}

type IConfig interface {
	GetSettings() IConfigSettings
	GetLogging() logger.ILogging
	GetAddress() IAddress
	GetConnection() string
}

type IConfigSettings interface {
	GetRetryNum() uint64
	GetPageOffset() uint64
	GetLanguage() language.ILanguage
}

type IAddress interface {
	GetExternal() string
}
