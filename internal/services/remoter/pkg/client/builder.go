package client

import (
	hls_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"
)

var (
	_ IBuilder = &sBuilder{}
)

type sBuilder struct {
	fPassword string
}

func NewBuilder(pPassword string) IBuilder {
	return &sBuilder{
		fPassword: pPassword,
	}
}

func (p *sBuilder) ExecCommand(pCmd ...string) *hls_settings.SCommandExecRequest {
	return &hls_settings.SCommandExecRequest{
		FPassword: p.fPassword,
		FCommand:  pCmd,
	}
}
