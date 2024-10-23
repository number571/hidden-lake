package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/network/anonymity"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleServicePubKeyAPI(pLogger logger.ILogger, pNode anonymity.INode) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(pkg_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		client := pNode.GetMessageQueue().GetClient()
		pubKey := client.GetPrivKey().GetPubKey()

		separated := pR.URL.Query().Get("separated")
		switch separated {
		case "true":
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, hls_settings.SPubKey{
				FKEMPKey: encoding.HexEncode(pubKey.GetKEMPubKey().ToBytes()),
				FDSAPKey: encoding.HexEncode(pubKey.GetDSAPubKey().ToBytes()),
			})
			return
		case "", "false":
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, pubKey.ToString())
			return
		default:
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_separeted"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: incorrect separated type")
			return
		}
	}
}
