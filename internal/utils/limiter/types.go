package limiter

type ILimitManager interface {
	Get(string) ILimiter
}

type ILimiter interface {
	Allow() bool
}
