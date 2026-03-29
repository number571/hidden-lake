package broker

import "context"

type IDataBroker interface {
	Register(string) error
	Produce(interface{})
	Consume(context.Context, string) (interface{}, error)
	CountSubscribers() uint64
}
