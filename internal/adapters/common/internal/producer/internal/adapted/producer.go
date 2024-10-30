package adapted

import (
	"bytes"
	"context"
	"errors"
	"net/http"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/adapters"
)

var (
	_ adapters.IAdaptedProducer = &sAdaptedProducer{}
)

type sAdaptedProducer struct {
	fServiceAddr string
}

func NewAdaptedProducer(pServiceAddr string) adapters.IAdaptedProducer {
	return &sAdaptedProducer{
		fServiceAddr: pServiceAddr,
	}
}

func (p *sAdaptedProducer) Produce(pCtx context.Context, pMsg net_message.IMessage) error {
	req, err := http.NewRequestWithContext(
		pCtx,
		http.MethodPost,
		p.fServiceAddr+"/push",
		bytes.NewBuffer([]byte(pMsg.ToString())),
	)
	if err != nil {
		return errors.Join(ErrBuildRequest, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Join(ErrBadRequest, err)
	}
	defer resp.Body.Close()
	if code := resp.StatusCode; code != http.StatusOK {
		return ErrBadStatusCode
	}
	return nil
}
