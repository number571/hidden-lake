package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/crypto/puzzle"
	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	hln_client "github.com/number571/hidden-lake/internal/applications/notifier/pkg/client"
	hln_settings "github.com/number571/hidden-lake/internal/applications/notifier/pkg/settings"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/alias"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleIncomingRedirectHTTP(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHLSClient hls_client.IClient,
) http.HandlerFunc {
	sett := pConfig.GetSettings()
	hlnClient := hln_client.NewClient(
		hln_client.NewSettings(&hln_client.SSettings{
			FDiffBits: sett.GetWorkSizeBits(),
			FParallel: sett.GetPowParallel(),
		}),
		hln_client.NewBuilder(),
		hln_client.NewRequester(pHLSClient),
	)

	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeOFF)

		logBuilder := http_logger.NewLogBuilder(hln_settings.GServiceName.Short(), pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		proof, _, saltBytes, bodyBytes, err := readRequestWithValidate(pR, sett.GetWorkSizeBits())
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("read_body"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: read body")
			return
		}

		friends, err := pHLSClient.GetFriends(pCtx)
		if err != nil || len(friends) < 2 {
			pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: get friends")
			return
		}

		fSender := asymmetric.LoadPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fSender == nil {
			pLogger.PushErro(logBuilder.WithMessage("load_pubkey"))
			_ = api.Response(pW, http.StatusForbidden, "failed: load public key")
			return
		}

		aliasName := alias.GetAliasByPubKey(friends, fSender)
		if aliasName == "" {
			pLogger.PushErro(logBuilder.WithMessage("find_alias"))
			_ = api.Response(pW, http.StatusForbidden, "failed: find alias")
			return
		}

		err = hlnClient.Redirect(pCtx, alias.GetAliasesList(friends), aliasName, proof, saltBytes, bodyBytes)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("redirect"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: redirect")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage("redirect"))
		_ = api.Response(pW, http.StatusOK, http_logger.CLogSuccess)
	}
}

func readRequestWithValidate(pR *http.Request, pWorkSizeBits uint64) (uint64, []byte, []byte, []byte, error) {
	bodyBytes, err := io.ReadAll(pR.Body)
	if err != nil {
		return 0, nil, nil, nil, err
	}

	proofBytes := encoding.HexDecode(pR.Header.Get(hln_settings.CHeaderPow))
	if proofBytes == nil {
		return 0, nil, nil, nil, errors.New("proof bytes is nil") // nolint: err113
	}

	saltBytes := encoding.HexDecode(pR.Header.Get(hln_settings.CHeaderSalt))
	if saltBytes == nil || len(saltBytes) != hln_client.CSaltSize {
		return 0, nil, nil, nil, errors.New("salt bytes is nil") // nolint: err113
	}

	proofArr := [encoding.CSizeUint64]byte{}
	copy(proofArr[:], proofBytes)
	proof := encoding.BytesToUint64(proofArr)

	hash := hashing.NewHasher(bytes.Join([][]byte{saltBytes, bodyBytes}, []byte{})).ToBytes()
	if !puzzle.NewPoWPuzzle(pWorkSizeBits).VerifyBytes(hash, proof) {
		return 0, nil, nil, nil, errors.New("salt bytes is nil") // nolint: err113
	}

	return proof, hash, saltBytes, bodyBytes, nil
}
