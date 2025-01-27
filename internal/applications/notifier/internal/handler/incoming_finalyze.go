package handler

import (
	"context"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/notifier/internal/database"
	"github.com/number571/hidden-lake/internal/utils/alias"
	"github.com/number571/hidden-lake/internal/utils/api"
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
	sett := pConfig.GetSettings()
	hlnClient := hln_client.NewClient(
		sett,
		hln_client.NewBuilder(),
		hln_client.NewRequester(pHLSClient),
	)

	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeOFF)

		logBuilder := http_logger.NewLogBuilder(hlm_settings.GServiceName.Short(), pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		proof, hash, saltBytes, bodyBytes, err := readRequestWithValidate(pR, sett.GetWorkSizeBits())
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("read_body"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: read body")
			return
		}

		myPubKey, err := pHLSClient.GetPubKey(pCtx)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		rel := database.NewRelation(myPubKey)
		hashExist, err := pDB.SetHash(rel, true, hash)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("set_hash"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: set hash")
			return
		}

		if !hashExist {
			dbMsg := database.NewMessage(true, bodyBytes)
			msg, err := msgdata.GetMessage(dbMsg.GetMessage(), dbMsg.GetTimestamp())
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("recv_message"))
				_ = api.Response(pW, http.StatusBadRequest, "failed: get message bytes")
				return
			}

			if err := pDB.Push(rel, dbMsg); err != nil {
				pLogger.PushErro(logBuilder.WithMessage("push_message"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: push message to database")
				return
			}

			pBroker.Produce("notifier", msg)
		}

		friends, err := pHLSClient.GetFriends(pCtx)
		if err != nil || len(friends) < 2 {
			pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: get friends")
			return
		}

		err = hlnClient.Finalyze(pCtx, alias.GetAliasesList(friends), proof, saltBytes, bodyBytes)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("finalyze"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: finalyze")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage("finalyze"))
		_ = api.Response(pW, http.StatusOK, http_logger.CLogSuccess)
	}
}
