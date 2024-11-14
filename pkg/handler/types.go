package handler

import (
	"context"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/hidden-lake/pkg/request"
	"github.com/number571/hidden-lake/pkg/response"
)

type IHandlerF func(
	context.Context,
	asymmetric.IPubKey,
	request.IRequest,
) (response.IResponse, error)
