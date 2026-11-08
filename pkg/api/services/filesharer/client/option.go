package client

type IClientOption func(*sClientOps)

type sClientOps struct {
	fPublic bool
}

func WithPublic(v bool) IClientOption {
	return func(o *sClientOps) {
		o.fPublic = v
	}
}
