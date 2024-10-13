package config

import "github.com/number571/hidden-lake/internal/helpers/encryptor/internal/config"

type IConfigSettings interface {
	config.IConfigSettings
}

type SConfigSettings struct {
	config.SConfigSettings
}
