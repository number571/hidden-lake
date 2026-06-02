package config

import "github.com/number571/hidden-lake/internal/adapters/meshtastic/pkg/app/config"

type IConfigSettings interface {
	config.IConfigSettings
}

type SConfigSettings struct {
	config.SConfigSettings
}
