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

	SetState(IMessageState)
	GetState() IMessageState

	GetUID() string
	GetBody() string

	GetIsNotify() bool
	GetTimestamp() string
}
