package process

type IPublishProcessManager interface {
	Unlock(string) bool
	TryLock(string) bool
}
