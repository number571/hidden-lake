package broker

type IDataBroker interface {
	Produce(interface{})
	Consume(string) <-chan interface{}
	CountSubscribers() uint64
}
