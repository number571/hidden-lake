package config

import "github.com/number571/hidden-lake/cmd/helpers/encryptor/internal/config"

type IConfigSettings interface {
	config.IConfigSettings
}

type SConfigSettings struct {
	config.SConfigSettings
}
