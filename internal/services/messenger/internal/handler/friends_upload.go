package handler

import (
	"context"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	hls_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/utils"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/webui"

	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
)

type sUploadFile struct {
	*sTemplate
	FAliasName    string
	FMessageLimit uint64
}

func FriendsUploadPage(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_messenger_settings.GetServiceName().Short(), pR)

		if pR.URL.Path != "/friends/upload" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		aliasName := pR.URL.Query().Get("alias_name")
		if aliasName == "" {
			ErrorPage(pLogger, pCfg, "get_alias_name", "alias name is nil")(pW, pR)
			return
		}

		msgLimit, err := utils.GetMessageLimit(pCtx, pHlsClient)
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_message_size", "get message size (limit)")(pW, pR)
			return
		}

		res := &sUploadFile{
			sTemplate:     getTemplate(pCfg),
			FAliasName:    aliasName,
			FMessageLimit: msgLimit,
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = webui.MustParseTemplate("index.html", "messenger/upload.html").Execute(pW, res)
	}
}
