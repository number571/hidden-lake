package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network"
	net_message "github.com/number571/go-peer/pkg/network/message"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/config"
	"github.com/number571/hidden-lake/internal/helpers/traffic/internal/storage"
	hlt_settings "github.com/number571/hidden-lake/internal/helpers/traffic/pkg/settings"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleMessageAPI(
	pCtx context.Context,
	pCfg config.IConfig,
	pStorage storage.IMessageStorage,
	pHTTPLogger, pAnonLogger logger.ILogger,
	pNode network.INode,
) http.HandlerFunc {
	tcpHandler := HandleServiceTCP(pCfg, pStorage, pAnonLogger)

	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlt_settings.CServiceName, pR)

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost {
			pHTTPLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		switch pR.Method {
		case http.MethodGet:
			query := pR.URL.Query()

			hash := encoding.HexDecode(query.Get("hash"))
			if len(hash) == 0 {
				pHTTPLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusTeapot, "failed: decode hash")
				return
			}

			msg, err := pStorage.Load(hash)
			if err != nil {
				pHTTPLogger.PushWarn(logBuilder.WithMessage("load_hash"))
				_ = api.Response(pW, http.StatusNotFound, "failed: load message")
				return
			}

			pHTTPLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, msg.ToString())
			return

		case http.MethodPost:
			msgStringAsBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pHTTPLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}

			netMsg, err := net_message.LoadMessage(
				pNode.GetSettings().GetConnSettings().GetMessageSettings(),
				string(msgStringAsBytes),
			)
			if err != nil {
				pHTTPLogger.PushWarn(logBuilder.WithMessage("decode_message"))
				_ = api.Response(pW, http.StatusTeapot, "failed: decode message")
				return
			}

			if netMsg.GetPayload().GetHead() != hls_settings.CNetworkMask {
				pHTTPLogger.PushWarn(logBuilder.WithMessage("network_mask"))
				_ = api.Response(pW, http.StatusLocked, "failed: network mask")
				return
			}

			if !pNode.GetCacheSetter().Set(netMsg.GetHash(), []byte{}) {
				pHTTPLogger.PushInfo(logBuilder.WithMessage("hash_already_exist"))
				_ = api.Response(pW, http.StatusAccepted, "accepted: hash already exist")
				return // hash of message already in queue
			}

			if err := tcpHandler(pCtx, pNode, nil, netMsg); err != nil {
				// internal logger
				_ = api.Response(pW, http.StatusBadRequest, "failed: handle message")
				return
			}

			pHTTPLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: handle message")
			return
		}
	}
}
