package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"

	anonymity "github.com/number571/go-peer/pkg/anonymity/qb"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/kernel/pkg/app/config"
	pkg_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/pkg/api/kernel/client/scheme"
)

func HandleConfigFriendsAPI(
	pWrapper config.IWrapper,
	pLogger logger.ILogger,
	pNode anonymity.INode,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.GetAppShortNameFMT(), pR)
		cfg := pWrapper.GetConfig()

		if pR.Method != http.MethodGet && pR.Method != http.MethodPost && pR.Method != http.MethodDelete {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		if pR.Method == http.MethodGet {
			friends := cfg.GetFriends()

			listFriends := make([][2]string, 0, len(friends))
			for name, pKey := range friends {
				listFriends = append(listFriends, [2]string{name, pKey.ToString()})
			}
			sort.Slice(listFriends, func(i, j int) bool {
				return listFriends[i][0] < listFriends[j][0]
			})

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, listFriends)
			return
		}

		friends := cfg.GetFriends()

		switch pR.Method {
		case http.MethodPost:
			var friend [2]string
			if err := json.NewDecoder(pR.Body).Decode(&friend); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}

			aliasName := friend[0]
			if aliasName == "" {
				pLogger.PushWarn(logBuilder.WithMessage("get_alias_name"))
				_ = api.Response(pW, http.StatusTeapot, "failed: load alias name")
				return
			}

			if _, ok := friends[aliasName]; ok {
				pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
				_ = api.Response(pW, http.StatusNotAcceptable, "failed: friend already exist")
				return
			}

			pKey := scheme.LoadParticipantKey(cfg.GetSettings().GetCryptoSchemeType(), friend[1])
			if pKey == nil {
				pLogger.PushWarn(logBuilder.WithMessage("decode_key"))
				_ = api.Response(pW, http.StatusBadRequest, "failed: load public key")
				return
			}

			friends[aliasName] = pKey
			if err := pWrapper.GetEditor().UpdateFriends(friends); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_friends"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: update friends")
				return
			}

			_ = pNode.GetKeysContainer().Add(pKey)

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: update friends")
			return

		case http.MethodDelete:
			aliasNameBytes, err := io.ReadAll(pR.Body)
			if err != nil {
				pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
				_ = api.Response(pW, http.StatusConflict, "failed: decode request")
				return
			}

			aliasName := string(aliasNameBytes)
			if aliasName == "" {
				pLogger.PushWarn(logBuilder.WithMessage("get_alias_name"))
				_ = api.Response(pW, http.StatusTeapot, "failed: load alias name")
				return
			}

			pKey, ok := friends[aliasName]
			if !ok {
				pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
				_ = api.Response(pW, http.StatusNotFound, "failed: friend does not exist")
				return
			}

			delete(friends, aliasName)

			if err := pWrapper.GetEditor().UpdateFriends(friends); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_friends"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: delete friend")
				return
			}

			_ = pNode.GetKeysContainer().Del(pKey.GetHasher().ToString())

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, "success: delete friend")
			return
		}
	}
}
