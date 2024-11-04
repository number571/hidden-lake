package config

import "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app/config"

type IConfigSettings interface {
	config.IConfigSettings
}

type SConfigSettings struct {
	config.SConfigSettings
}
