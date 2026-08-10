package dto

import "github.com/number571/go-peer/pkg/types"

type IMessage interface {
	types.IConverter

	IsIncoming() bool
	GetTimestamp() string
	GetMessage() string
}
