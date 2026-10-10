package process

type INotifyProcessManager interface {
	Unlock(string) bool
	TryLock(string) bool
}
