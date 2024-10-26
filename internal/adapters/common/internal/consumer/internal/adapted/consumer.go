// nolint: err113
package adapted

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/adapters"
)

var (
	_ adapters.IAdaptedConsumer = &sAdaptedConsumer{}
)

type sAdaptedConsumer struct {
	fSettings    net_message.ISettings
	fServiceAddr string
	fStates      [2]sState // 0 = got, 1 = cur
}

type sState struct {
	fIter uint64
	fLast uint64
}

func NewAdaptedConsumer(
	pSettings net_message.ISettings,
	pServiceAddr string,
) adapters.IAdaptedConsumer {
	return &sAdaptedConsumer{
		fSettings:    pSettings,
		fServiceAddr: pServiceAddr,
		fStates:      [2]sState{},
	}
}

func (p *sAdaptedConsumer) Consume(pCtx context.Context) (net_message.IMessage, error) {
	if p.fStates[0].fIter == p.fStates[1].fIter && p.fStates[0].fLast == p.fStates[1].fLast {
		iter, last, err := p.getIterAndLast(pCtx)
		if err != nil {
			return nil, err
		}
		if iter == p.fStates[0].fIter && last == p.fStates[0].fLast {
			return nil, nil
		}
		p.fStates[0].fIter = iter
		p.fStates[0].fLast = last
	}

	if p.fStates[1].fIter != p.fStates[0].fIter {
		p.fStates[1].fIter = p.fStates[0].fIter
		p.fStates[1].fLast = 0
	}

	msg, err := p.loadMessage(pCtx)
	if err != nil {
		p.fStates[0].fIter, p.fStates[0].fLast = 0, 0
		p.fStates[1].fIter, p.fStates[1].fLast = 0, 0
		return nil, err
	}

	p.fStates[1].fLast++
	return msg, nil
}

func (p *sAdaptedConsumer) getIterAndLast(pCtx context.Context) (uint64, uint64, error) {
	req, err := http.NewRequestWithContext(pCtx, http.MethodGet, p.fServiceAddr+"/last", nil)
	if err != nil {
		return 0, 0, err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer rsp.Body.Close()
	if code := rsp.StatusCode; code != http.StatusOK {
		return 0, 0, fmt.Errorf("status code: %d", code)
	}
	res, err := io.ReadAll(rsp.Body)
	if err != nil {
		return 0, 0, err
	}
	splited := strings.Split(string(res), ".")
	if len(splited) != 2 {
		return 0, 0, errors.New("failed splited uints")
	}
	iter, err := strconv.ParseUint(splited[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed iter: %s", err.Error())
	}
	last, err := strconv.ParseUint(splited[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed last: %s", err.Error())
	}
	return iter, last, nil
}

func (p *sAdaptedConsumer) loadMessage(pCtx context.Context) (net_message.IMessage, error) {
	req, err := http.NewRequestWithContext(
		pCtx,
		http.MethodGet,
		fmt.Sprintf(p.fServiceAddr+"/load?id=%d", p.fStates[1].fLast),
		nil,
	)
	if err != nil {
		return nil, err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if code := rsp.StatusCode; code != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", code)
	}
	res, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	return net_message.LoadMessage(p.fSettings, string(res))
}
