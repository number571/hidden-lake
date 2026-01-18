package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/network"
	"github.com/number571/hidden-lake/pkg/network/request"
)

func HandleNetworkRequestAPI(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pNode network.IHiddenLakeNode,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetAppShortNameFMT(), pR)

		var vRequest = &request.SRequest{}

		if pR.Method != http.MethodPost && pR.Method != http.MethodPut {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if err := json.NewDecoder(pR.Body).Decode(vRequest); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: decode request")
			return
		}

		friends := pConfig.GetFriends()
		pubKey, ok := friends[pR.URL.Query().Get("friend")]
		if !ok {
			pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: load public key")
			return
		}

		switch pR.Method {
		case http.MethodPut:
			if err := pNode.SendRequest(pCtx, pubKey, vRequest); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("send_payload"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: send payload")
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: send")
			return

		case http.MethodPost:
			resp, err := pNode.FetchRequest(pCtx, pubKey, vRequest)
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
