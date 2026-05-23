package handler

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/scheme/layer2"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

type IHandlerF func(
	context.Context,
	layer2.IParticipantKey,
	request.IRequest,
) (response.IResponse, error)
