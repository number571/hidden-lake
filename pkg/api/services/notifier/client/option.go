package client

type IClientOption func(*sClientOps)

type sClientOps struct {
	fWaitResponse bool
}

func WithWaitResponse(v bool) IClientOption {
	return func(o *sClientOps) {
		o.fWaitResponse = v
	}
}
