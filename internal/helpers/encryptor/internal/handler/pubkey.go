package handler

import (
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"

	hle_settings "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
)

func HandleServicePubKeyAPI(pLogger logger.ILogger, pPubKey asymmetric.IPubKey) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hle_settings.CServiceName, pR)

		if pR.Method != http.MethodGet {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		separated := pR.URL.Query().Get("separated")
		switch separated {
		case "true":
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, [2]string{
				encoding.HexEncode(pPubKey.GetKEMPubKey().ToBytes()),
				encoding.HexEncode(pPubKey.GetDSAPubKey().ToBytes()),
			})
			return
		case "", "false":
			pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
			_ = api.Response(pW, http.StatusOK, pPubKey.ToString())
			return
		default:
			pLogger.PushWarn(logBuilder.WithMessage("incorrect_separeted"))
			_ = api.Response(pW, http.StatusBadRequest, "failed: incorrect separated type")
			return
		}
	}
}
