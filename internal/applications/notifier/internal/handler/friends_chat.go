package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/notifier/internal/database"
	internal_utils "github.com/number571/hidden-lake/internal/applications/notifier/internal/utils"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	hln_client "github.com/number571/hidden-lake/internal/applications/notifier/pkg/client"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/notifier/pkg/settings"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	"github.com/number571/hidden-lake/internal/utils/alias"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/msgdata"
	"github.com/number571/hidden-lake/internal/webui"
)

type sChatMessage struct {
	FIsIncoming bool
	msgdata.SMessage
}

type sFriendsChat struct {
	*sTemplate
	FMessages []sChatMessage
}

func FriendsChatPage(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pDB database.IKVDatabase,
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.GServiceName.Short(), pR)

		if pR.URL.Path != "/friends/chat" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		if err := pR.ParseForm(); err != nil {
			ErrorPage(pLogger, pCfg, "parse_form", "parse form")(pW, pR)
			return
		}

		// default max value = 16MiB
		if err := pR.ParseMultipartForm(16 << 20); err != nil && !errors.Is(err, http.ErrNotMultipart) {
			ErrorPage(pLogger, pCfg, "parse_multipart_form", "parse multipart form")(pW, pR)
			return
		}

		myPubKey, err := pHlsClient.GetPubKey(pCtx)
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_public_key", "read public key")(pW, pR)
			return
		}

		friends, err := pHlsClient.GetFriends(pCtx)
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_friends", "read friends list")(pW, pR)
			return
		}

		rel := database.NewRelation(myPubKey)

		switch pR.FormValue("method") {
		case http.MethodPost, http.MethodPut:
			msgBytes, err := msgdata.GetMessageBytes(pR)
			if err != nil || msgBytes == nil {
				ErrorPage(pLogger, pCfg, "get_message", "get message bytes")(pW, pR)
				return
			}

			hash, err := pushMessage(pCtx, pCfg, pHlsClient, alias.GetAliasesList(friends), msgBytes)
			if err != nil {
				ErrorPage(pLogger, pCfg, "send_message", "push message to network")(pW, pR)
				return
			}

			if _, err := pDB.SetHash(rel, false, hash); err != nil {
				ErrorPage(pLogger, pCfg, "set_hash", "set hash of message to database")(pW, pR)
				return
			}

			dbMsg := database.NewMessage(false, msgBytes)
			if err := pDB.Push(rel, dbMsg); err != nil {
				ErrorPage(pLogger, pCfg, "push_message", "add message to database")(pW, pR)
				return
			}

			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogRedirect))
			http.Redirect(pW, pR, "/friends/chat", http.StatusSeeOther)
			return
		}

		start := uint64(0)
		size := pDB.Size(rel)

		messagesCap := pCfg.GetSettings().GetMessagesCapacity()
		if size > messagesCap {
			start = size - messagesCap
		}

		dbMsgs, err := pDB.Load(rel, start, size)
		if err != nil {
			ErrorPage(pLogger, pCfg, "read_database", "read database")(pW, pR)
			return
		}

		res := &sFriendsChat{
			sTemplate: getTemplate(pCfg),
			FMessages: func() []sChatMessage {
				msgs := make([]sChatMessage, 0, len(dbMsgs))
				for _, dbMsg := range dbMsgs {
					msg, err := msgdata.GetMessage(dbMsg.GetMessage(), dbMsg.GetTimestamp())
					if err != nil {
						panic(err)
					}
					msgs = append(msgs, sChatMessage{
						FIsIncoming: dbMsg.IsIncoming(),
						SMessage:    msg,
					})
				}
				return msgs
			}(),
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = webui.MustParseTemplate("index.html", "notifier/chat.html").Execute(pW, res)
	}
}

func pushMessage(
	pCtx context.Context,
	pCfg config.IConfig,
	pClient hls_client.IClient,
	pFriends []string,
	pMsgBytes []byte,
) ([]byte, error) {
	msgLimit, err := internal_utils.GetMessageLimit(pCtx, pClient)
	if err != nil {
		return nil, errors.Join(ErrGetMessageLimit, err)
	}

	if uint64(len(pMsgBytes)) > msgLimit {
		return nil, ErrLenMessageGtLimit
	}

	sett := pCfg.GetSettings()
	hlnClient := hln_client.NewClient(
		hln_client.NewSettings(&hln_client.SSettings{
			FDiffBits: sett.GetWorkSizeBits(),
			FParallel: sett.GetPowParallel(),
		}),
		hln_client.NewBuilder(),
		hln_client.NewRequester(pClient),
	)

	hash, err := hlnClient.Initialize(pCtx, pFriends, pMsgBytes)
	if err != nil {
		return nil, errors.Join(ErrPushMessage, err)
	}

	return hash, nil
}
