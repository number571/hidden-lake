package client

import (
	"context"
)

type IClient interface {
	GetIndex(context.Context) error

	PingFriend(context.Context, string) error
}

type IRequester interface {
	GetIndex(context.Context) error

	PingFriend(context.Context, string) error
}
