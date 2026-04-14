package config

import "github.com/number571/hidden-lake/internal/adapters/https/pkg/app/config"

type IConfigSettings interface {
	config.IConfigSettings
}

type SConfigSettings struct {
	config.SConfigSettings
}
