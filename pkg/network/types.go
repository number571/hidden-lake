package network

import (
	"time"

	gopeer_logger "github.com/number571/go-peer/pkg/logger"
	gopeer_message "github.com/number571/go-peer/pkg/network/message"
)

type ISettings interface {
	ISubSettings
	GetMessageSettings() gopeer_message.ISettings
	GetMessageSizeBytes() uint64
	GetQueuePeriod() time.Duration
	GetFetchTimeout() time.Duration
}

type ISubSettings interface {
	GetLogger() gopeer_logger.ILogger
	GetParallel() uint64
	GetTCPAddress() string
}
