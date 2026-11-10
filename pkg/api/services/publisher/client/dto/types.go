package dto

import "github.com/number571/go-peer/pkg/types"

type IRemoteType uint8

const (
	CRemoteDisabled IRemoteType = 0 + iota
	CRemoteEnabled
	CRemoteEnabledWithoutWaitResponse
	cEndRemoteUndefined
)

type IMessageState uint8

const (
	cBeginFlagState IMessageState = iota
	CIncomingState
	CSendingState
	CDeliveredState
	CPublishedState
	cEndFlagState
)

type IMessage interface {
	types.IConverter

	GetHead() IMessageHead
	GetBody() IMessageBody
}

type IMessageHead interface {
	types.IConverter

	GetState() IMessageState
	GetUID() string
}

type IMessageBody interface {
	types.IConverter

	GetData() []byte
	GetTimestamp() string
}
