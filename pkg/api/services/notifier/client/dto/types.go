package dto

import "github.com/number571/go-peer/pkg/types"

type IMessageState uint8

const (
	cBeginFlagState IMessageState = iota
	CIncomingState
	CSendingState
	CReceivedState
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

	GetData() string
	GetTimestamp() string
}
