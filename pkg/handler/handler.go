package handler

import (
	"context"
	"errors"

	"github.com/number571/go-peer/pkg/anonymity"
	anon_logger "github.com/number571/go-peer/pkg/anonymity/logger"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	internal_anon_logger "github.com/number571/hidden-lake/internal/utils/logger/anon"
	"github.com/number571/hidden-lake/pkg/request"
)

func RequestHandler(pHandleF IHandlerF) anonymity.IHandlerF {
	return func(
		pCtx context.Context,
		pNode anonymity.INode,
		pSender asymmetric.IPubKey,
		pReqBytes []byte,
	) ([]byte, error) {
		logger := pNode.GetLogger()
		logBuilder := anon_logger.NewLogBuilder(pNode.GetSettings().GetServiceName())

		// enrich logger
		logBuilder.
			WithSize(len(pReqBytes)).
			WithPubKey(pSender)

		// load request from message's body
		loadReq, err := request.LoadRequest(pReqBytes)
		if err != nil {
			logger.PushErro(logBuilder.WithType(internal_anon_logger.CLogErroLoadRequestType))
			return nil, errors.Join(ErrLoadRequest, err)
		}

		// handle request
		rsp, err := pHandleF(pCtx, pSender, loadReq)
		if err != nil {
			logger.PushWarn(logBuilder.WithType(internal_anon_logger.CLogWarnRequestHandle))
			return nil, ErrUndefinedService
		}

		// no need response
		if rsp == nil {
			logger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogBaseResponseModeFromService))
			return nil, nil
		}

		// send response
		logger.PushInfo(logBuilder.WithType(internal_anon_logger.CLogInfoResponseFromService))
		return rsp.ToBytes(), nil
	}
}
