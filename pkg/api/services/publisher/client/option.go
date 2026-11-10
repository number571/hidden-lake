package client

import "github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"

type IClientOption func(*sClientOps)

type sClientOps struct {
	fRemoteType dto.IRemoteType
	fPublic     bool
}

func WithRemoteType(v dto.IRemoteType) IClientOption {
	return func(o *sClientOps) {
		o.fRemoteType = v
	}
}

func WithPublic(v bool) IClientOption {
	return func(o *sClientOps) {
		o.fPublic = v
	}
}
