package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	hls_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/database"
	internal_utils "github.com/number571/hidden-lake/internal/services/messenger/internal/utils"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	hls_messenger_client "github.com/number571/hidden-lake/internal/services/messenger/pkg/client"
	hls_pinger_client "github.com/number571/hidden-lake/internal/services/pinger/pkg/client"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/msgdata"
	"github.com/number571/hidden-lake/internal/webui"

	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
)

type sChatMessage struct {
	FIsIncoming bool
	msgdata.SMessage
}

type sChatAddress struct {
	FAliasName  string
	FPublicKey  string
	FPubKeyHash string
}

type sFriendsChat struct {
	*sTemplate
	FPingState int
	FAddress   sChatAddress
	FMessages  []sChatMessage
}

func FriendsChatPage(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pDB database.IKVDatabase,
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_messenger_settings.GetShortAppName(), pR)

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

		aliasName := pR.URL.Query().Get("alias_name")
		if aliasName == "" {
			ErrorPage(pLogger, pCfg, "get_alias_name", "alias name is nil")(pW, pR)
			return
		}

		myPubKey, err := pHlsClient.GetPubKey(pCtx)
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_public_key", "read public key")(pW, pR)
			return
		}

		recvPubKey, err := getReceiverPubKey(pCtx, pHlsClient, aliasName)
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_receiver", "get receiver by public key")(pW, pR)
			return
		}

		pingState := 0
		rel := database.NewRelation(myPubKey, recvPubKey)

		switch pR.FormValue("method") {
		case http.MethodPost, http.MethodPut:
			msgBytes, err := msgdata.GetMessageBytes(pR)
			if err != nil {
				ErrorPage(pLogger, pCfg, "get_message", "get message bytes")(pW, pR)
				return
			}

			if msgBytes != nil {
				if err := pushMessage(pCtx, pHlsClient, aliasName, msgBytes); err != nil {
					ErrorPage(pLogger, pCfg, "send_message", "push message to network")(pW, pR)
					return
				}
				dbMsg := database.NewMessage(false, msgBytes)
				if err := pDB.Push(rel, dbMsg); err != nil {
					ErrorPage(pLogger, pCfg, "push_message", "add message to database")(pW, pR)
					return
				}
				pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogRedirect))
				http.Redirect(pW, pR, "/friends/chat?alias_name="+aliasName, http.StatusSeeOther)
				return
			}

			hlpClient := hls_pinger_client.NewClient(
				hls_pinger_client.NewBuilder(),
				hls_pinger_client.NewRequester(pHlsClient),
			)

			pingState = 1
			if err := hlpClient.Ping(pCtx, aliasName); err != nil {
				pingState = -1
			}
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
			sTemplate:  getTemplate(pCfg),
			FPingState: pingState,
			FAddress: sChatAddress{
				FAliasName:  aliasName,
				FPublicKey:  recvPubKey.ToString(),
				FPubKeyHash: recvPubKey.GetHasher().ToString(),
			},
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
		_ = webui.MustParseTemplate("index.html", "messenger/chat.html").Execute(pW, res)
	}
}

func pushMessage(
	pCtx context.Context,
	pClient hls_client.IClient,
	pAliasName string,
	pMsgBytes []byte,
) error {
	msgLimit, err := internal_utils.GetMessageLimit(pCtx, pClient)
	if err != nil {
		return errors.Join(ErrGetMessageLimit, err)
	}

	if uint64(len(pMsgBytes)) > msgLimit {
		return ErrLenMessageGtLimit
	}

	hlmClient := hls_messenger_client.NewClient(
		hls_messenger_client.NewBuilder(),
		hls_messenger_client.NewRequester(pClient),
	)

	if err := hlmClient.PushMessage(pCtx, pAliasName, pMsgBytes); err != nil {
		return errors.Join(ErrPushMessage, err)
	}

	return nil
}

func getReceiverPubKey(
	pCtx context.Context,
	client hls_client.IClient,
	aliasName string,
) (asymmetric.IPubKey, error) {
	friends, err := client.GetFriends(pCtx)
	if err != nil {
		return nil, errors.Join(ErrGetFriends, err)
	}

	friendPubKey, ok := friends[aliasName]
	if !ok {
		return nil, ErrUndefinedPublicKey
	}

	return friendPubKey, nil
}
