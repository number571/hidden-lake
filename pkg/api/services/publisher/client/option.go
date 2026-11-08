package client

type IClientOption func(*sClientOps)

type sClientOps struct {
	fRemote bool
	fPublic bool
}

func WithRemote(v bool) IClientOption {
	return func(o *sClientOps) {
		o.fRemote = v
	}
}

func WithPublic(v bool) IClientOption {
	return func(o *sClientOps) {
		o.fPublic = v
	}
}
