package network

import (
	"context"
	"crypto/ed25519"

	"github.com/number571/go-peer/pkg/types"
)

type ICallbackFunc func(ed25519.PublicKey, []byte, string)

type IHiddenLakeChatNode interface {
	types.IRunner
	GetMessageLimitSize() uint64
	SendMessage(context.Context, string) error
}
