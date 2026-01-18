package handler

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/hidden-lake/pkg/network/request"
	"github.com/number571/hidden-lake/pkg/network/response"
)

type IHandlerF func(
	context.Context,
	asymmetric.IPubKey,
	request.IRequest,
) (response.IResponse, error)
