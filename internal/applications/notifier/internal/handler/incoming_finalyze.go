package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/hidden-lake/internal/applications/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/utils/alias"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/layer3"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/msgdata"

	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	hln_client "github.com/number571/hidden-lake/internal/applications/notifier/pkg/client"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/notifier/pkg/settings"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func HandleIncomingFinalyzeHTTP(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pDB database.IKVDatabase,
	pBroker msgdata.IMessageBroker,
	pHLSClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeOFF)

		logBuilder := http_logger.NewLogBuilder(hlm_settings.GServiceName.Short(), pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		msgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusConflict, "failed: decode body")
			return
		}

		rawMsg, err := layer1.LoadMessage(
			layer1.NewSettings(&layer1.SSettings{
				FWorkSizeBits: pConfig.GetSettings().GetWorkSizeBits(),
			}),
			msgBytes,
		)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_message"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: decode message")
			return
		}

		rawBodyBytes, err := layer3.ExtractMessageBody(rawMsg)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("extract_raw_message"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: extract raw message")
			return
		}

		myPubKey, err := pHLSClient.GetPubKey(pCtx)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		hashExist, err := pDB.SetHash(myPubKey, true, rawMsg.GetHash())
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("set_hash"))
			_ = api.Response(pW, http.StatusInternalServerError, "failed: set hash")
			return
		}

		decMsg, channelKey := tryDecryptMessage(pConfig, rawBodyBytes)
		if !hashExist && channelKey != "" {
			if _, err := pDB.SetHash(myPubKey, true, decMsg.GetHash()); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("set_hash"))
				_ = api.Response(pW, http.StatusNotAcceptable, "failed: set hash")
				return
			}

			bodyBytes, err := layer3.ExtractMessageBody(decMsg)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("extract_dec_message"))
				_ = api.Response(pW, http.StatusBadRequest, "failed: extract dec message")
				return
			}

			dbMsg := database.NewMessage(true, bodyBytes)
			msg, err := msgdata.GetMessage(dbMsg.GetMessage(), dbMsg.GetTimestamp())
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("recv_message"))
				_ = api.Response(pW, http.StatusBadRequest, "failed: get message bytes")
				return
			}

			rel := database.NewRelation(myPubKey, channelKey)
			if err := pDB.Push(rel, dbMsg); err != nil {
				pLogger.PushErro(logBuilder.WithMessage("push_message"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
				return
			}

			pBroker.Produce(channelKey, msg)
		}

		friends, err := pHLSClient.GetFriends(pCtx)
		if err != nil || len(friends) < 2 {
			pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: get friends")
			return
		}

		hlnClient := hln_client.NewClient(
			hln_client.NewBuilder(),
			hln_client.NewRequester(pHLSClient),
		)
		if err := hlnClient.Finalyze(pCtx, alias.GetAliasesList(friends), rawMsg); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("finalyze"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: finalyze")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage("finalyze"))
		_ = api.Response(pW, http.StatusOK, http_logger.CLogSuccess)
	}
}

func tryDecryptMessage(pConfig config.IConfig, pBody []byte) (layer1.IMessage, string) {
	for _, key := range pConfig.GetChannels() {
		// try decrypt message
		decMsg, err := layer1.LoadMessage(
			layer1.NewSettings(&layer1.SSettings{
				FNetworkKey: key,
			}),
			pBody,
		)
		if err != nil {
			continue
		}
		return decMsg, key
	}
	return nil, ""
}
