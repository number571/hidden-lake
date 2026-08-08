package network

import (
	"context"
	"errors"
	"sync"

	"github.com/number571/go-peer/pkg/anonymity/qb/queue"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer1"
	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/go-peer/pkg/state"
)

var (
	_ queue.IQBProblemProcessor = &sSilentQueueProcessor{}
)

type sSilentQueueProcessor struct {
	fState state.IState

	fSettings queue.ISettings
	fScheme   layer2.IScheme

	fInQueue  chan []byte
	fOutQueue chan layer1.IMessage
}

func newSilentQueueProcessor(pSettings queue.ISettings, pScheme layer2.IScheme) queue.IQBProblemProcessor {
	queuePoolCap := pSettings.GetQueuePoolCap()[0]
	return &sSilentQueueProcessor{
		fSettings: pSettings,
		fScheme:   pScheme,
		fInQueue:  make(chan []byte, queuePoolCap),
		fOutQueue: make(chan layer1.IMessage, queuePoolCap),
	}
}

func (p *sSilentQueueProcessor) Run(pCtx context.Context) error {
	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	if err := p.fState.Enable(nil); err != nil {
		return errors.Join(ErrRunning, err)
	}
	defer func() { _ = p.fState.Disable(nil) }()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go p.runMainPoolFiller(ctx, cancel, &wg)

	wg.Wait()
	return ctx.Err()
}

func (p *sSilentQueueProcessor) GetSettings() queue.ISettings {
	return p.fSettings
}

func (p *sSilentQueueProcessor) GetScheme() layer2.IScheme {
	return p.fScheme
}

func (p *sSilentQueueProcessor) EnqueueMessage(pKey layer2.IParticipantKey, pBytes []byte) error {
	rawMsg, err := p.fScheme.EncryptMessage(pKey, pBytes)
	if err != nil {
		return errors.Join(ErrEncryptMessage, err)
	}
	p.fInQueue <- rawMsg
	return nil
}

func (p *sSilentQueueProcessor) DequeueMessage(pCtx context.Context) layer1.IMessage {
	for {
		select {
		case <-pCtx.Done():
			return nil
		case msg := <-p.fOutQueue:
			return msg
		}
	}
}

func (p *sSilentQueueProcessor) runMainPoolFiller(pCtx context.Context, pCancel func(), pWG *sync.WaitGroup) {
	defer func() {
		pWG.Done()
		pCancel()
	}()
	for {
		select {
		case <-pCtx.Done():
			return
		case rawMsg := <-p.fInQueue:
			if err := p.pushMessage(pCtx, rawMsg); err != nil {
				return
			}
		}
	}
}

func (p *sSilentQueueProcessor) pushMessage(pCtx context.Context, pRawMsg []byte) error {
	chNetMsg := make(chan layer1.IMessage)
	go func() {
		chNetMsg <- layer1.NewMessage(
			p.fSettings.GetMessageConstructSettings(),
			pRawMsg,
		)
	}()
	select {
	case <-pCtx.Done():
		return pCtx.Err()
	case netMsg := <-chNetMsg:
		p.fOutQueue <- netMsg
		return nil
	}
}
