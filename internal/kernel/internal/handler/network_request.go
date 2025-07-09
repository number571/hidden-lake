package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/request"
)

const (
	cErrorNone = iota
	cErrorGetFriends
	cErrorDecodeData
)

func HandleNetworkRequestAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pNode network.IHiddenLakeNode,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetServiceName().Short(), pR)

		var vRequest pkg_settings.SRequest

		if pR.Method != http.MethodPost && pR.Method != http.MethodPut {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(&vRequest); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		pubKey, req, errCode := unwrapRequest(pConfig, vRequest)
		switch errCode {
		case cErrorNone:
			// pass
		case cErrorGetFriends:
			pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: load public key")
			return
		case cErrorDecodeData:
			pLogger.PushWarn(logBuilder.WithMessage("decode_data"))
			_ = api.Response(pW, http.StatusTeapot, "failed: decode hex format data")
			return
		default:
			panic("undefined error code")
		}

		switch pR.Method {
		case http.MethodPut:
			if err := pNode.SendRequest(pCtx, pubKey, req); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("send_payload"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: send payload")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: send")
			return

		case http.MethodPost:
			resp, err := pNode.FetchRequest(pCtx, pubKey, req)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("fetch_payload"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: fetch payload")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, resp)
			return
		}
	}
}

func unwrapRequest(
	pConfig config.IConfig,
	pRequest pkg_settings.SRequest,
) (asymmetric.IPubKey, request.IRequest, int) {
	friends := pConfig.GetFriends()

	pubKey, ok := friends[pRequest.FReceiver]
	if !ok {
		return nil, nil, cErrorGetFriends
	}

	req := pRequest.FReqData
	if req == nil {
		return nil, nil, cErrorDecodeData
	}

	return pubKey, req, cErrorNone
}
