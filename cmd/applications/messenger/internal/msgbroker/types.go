package msgbroker

import "github.com/number571/hidden-lake/cmd/applications/messenger/internal/utils"

type IMessageBroker interface {
	Produce(string, utils.SMessage)
	Consume(string) (utils.SMessage, bool)
}
