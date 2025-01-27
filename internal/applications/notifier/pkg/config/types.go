package config

import "github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"

type IConfigSettings interface {
	config.IConfigSettings
}

type SConfigSettings struct {
	config.SConfigSettings
}
