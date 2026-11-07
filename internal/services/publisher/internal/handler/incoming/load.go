package incoming

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/publisher/internal/database"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/broker"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	message_dto "github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/services/publisher/pkg/app/config"
	hls_publisher_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
)

func HandleIncomingLoadHTTP(
	pCtx context.Context,
	pLogger logger.ILogger,
	pConfig config.IConfig,
	pDB database.IKVDatabase,
	pBroker broker.IDataBroker,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hlk_settings.CHeaderResponseMode, hlk_settings.CHeaderResponseModeON)

		logBuilder := http_logger.NewLogBuilder(hls_publisher_settings.GetAppShortNameFMT(), pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "method", api.WithoutContentType())
			return
		}

		aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName)
		if aliasName == "" {
			pLogger.PushErro(logBuilder.WithMessage("alias_name_is_null"))
			_ = api.Response(pW, http.StatusBadGateway, "alias", api.WithoutContentType())
			return
		}

		queryParams := pR.URL.Query()
		var (
			msgNID = queryParams.Get("nid")
			msgUID = queryParams.Get("uid")
		)

		if msgNID != "" && msgUID != "" {
			pLogger.PushWarn(logBuilder.WithMessage("has_many_ids"))
			_ = api.Response(pW, http.StatusBadRequest, "manyids")
			return
		}

		var index uint64
		switch {
		case msgNID != "":
			var err error
			index, err = strconv.ParseUint(msgNID, 10, 64)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("msg_index"))
				_ = api.Response(pW, http.StatusBadRequest, "nid", api.WithoutContentType())
				return
			}
		case msgUID != "":
			if !message_dto.MessageUIDIsValid(msgUID) {
				pLogger.PushWarn(logBuilder.WithMessage("msg_uid"))
				_ = api.Response(pW, http.StatusBadRequest, "uid", api.WithoutContentType())
				return
			}
		default:
			pLogger.PushWarn(logBuilder.WithMessage("id_undefined"))
			_ = api.Response(pW, http.StatusBadRequest, "idundef", api.WithoutContentType())
			return
		}

		var (
			msg message_dto.IMessage
			err error
		)
		if msgUID != "" {
			msg, err = pDB.LoadByUID(aliasName, msgUID)
		} else {
			msg, err = pDB.LoadByEnum(aliasName, index)
		}
		if err == nil {
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, msg.ToBytes(), api.WithoutContentType())
			return
		}
		if errors.Is(err, database.ErrNotFoundMessage) {
			pLogger.PushWarn(logBuilder.WithMessage("message_not_found"))
			_ = api.Response(pW, http.StatusNotFound, "notfound", api.WithoutContentType())
			return
		}

		pLogger.PushErro(logBuilder.WithMessage("load_message"))
		_ = api.Response(pW, http.StatusInternalServerError, "getmsg", api.WithoutContentType())
	}
}
