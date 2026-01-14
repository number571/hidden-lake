package handler

import (
	"context"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/utils"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleLoadMessagesAPI(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pHlkClient hlk_client.IClient,
	pDatabase database.IKVDatabase,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		fPubKey, err := utils.GetFriendPubKeyByAliasName(pCtx, pHlkClient, pR.URL.Query().Get("friend"))
		if err != nil {
			pLogger.PushErro(logBuilder.WithMessage("load_pubkey"))
			_ = api.Response(pW, http.StatusForbidden, "failed: load public key")
			return
		}

		myPubKey, err := pHlkClient.GetPubKey(pCtx)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		rel := database.NewRelation(myPubKey, fPubKey)

		// TODO:
		// pR.URL.Query().Get("page")
		// pR.URL.Query().Get("offset")

		start := uint64(0)
		size := pDatabase.Size(rel)

		messagesCap := pConfig.GetSettings().GetMessagesCapacity()
		if size > messagesCap {
			start = size - messagesCap
		}

		dbMsgs, err := pDatabase.Load(rel, start, size)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, dbMsgs)
	}
}
