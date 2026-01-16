package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/pubkey"
)

func HandleChatHistoryLoadAPI(
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

		queryParams := pR.URL.Query()
		fPubKey, err := pubkey.GetFriendPubKeyByAliasName(pCtx, pHlkClient, queryParams.Get("friend"))
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
		size := pDatabase.Size(rel)

		count := pConfig.GetSettings().GetMessagesCapacity()
		if x := queryParams.Get("count"); x != "" {
			var err error
			count, err = strconv.ParseUint(x, 10, 64)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("parse_count"))
				_ = api.Response(pW, http.StatusBadGateway, "failed: parse count")
				return
			}
		}

		start := uint64(0)
		if x := queryParams.Get("start"); x != "" {
			var err error
			start, err = strconv.ParseUint(x, 10, 64)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("parse_start"))
				_ = api.Response(pW, http.StatusBadGateway, "failed: parse start")
				return
			}
		}

		// ASC select used by default
		if x := queryParams.Get("select"); x == "desc" {
			if size < start {
				start = 0
			} else {
				start = (size - start)
			}
		}

		dbMsgs, err := pDatabase.Load(rel, start, count)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("get_public_key"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: get public key from service")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = api.Response(pW, http.StatusOK, dbMsgs)
	}
}
