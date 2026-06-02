package meshtastic

import (
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/pkg/network/adapters"
)

type IMeshtasticAdapter interface {
	adapters.IRunnerAdapter

	WithLogger(string, logger.ILogger) IMeshtasticAdapter
}

type ISettings interface {
	ISrvSettings

	GetAdapterSettings() adapters.ISettings
}

type ISrvSettings interface {
	GetPath() string
	GetAddress() string
	GetChannel() uint8
	GetDevPath() string
}
