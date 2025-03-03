package anon

import anon_logger "github.com/number571/go-peer/pkg/anonymity/logger"

var gLogMap = map[anon_logger.ILogType]string{
	0: "", // invalid log

	// default
	anon_logger.CLogBaseBroadcast:         "BRDCS",
	anon_logger.CLogBaseEnqueueRequest:    "ENQRQ",
	anon_logger.CLogBaseEnqueueResponse:   "ENQRS",
	anon_logger.CLogBaseGetResponse:       "GETRS",
	anon_logger.CLogInfoExist:             "EXIST",
	anon_logger.CLogInfoUndecryptable:     "UNDEC",
	anon_logger.CLogInfoWithoutResponse:   "WTHRS",
	anon_logger.CLogWarnMessageNull:       "MNULL",
	anon_logger.CLogWarnPayloadNull:       "PNULL",
	anon_logger.CLogWarnUnknownRoute:      "UNKRT",
	anon_logger.CLogWarnIncorrectResponse: "ICRSP",
	anon_logger.CLogErroDatabaseGet:       "DBGET",
	anon_logger.CLogErroDatabaseSet:       "DBSET",

	// extend
	CLogBaseResponseModeFromService: "RSPMD",
	CLogBaseSendNetworkMessage:      "SNMSG",
	CLogInfoResponseFromService:     "RSPSR",
	CLogInfoRecvNetworkMessage:      "RNMSG",
	CLogWarnRequestToService:        "RQTSR",
	CLogWarnUndefinedService:        "UNDSR",
	CLogWarnInvalidRequestMethod:    "IRMTH",
	CLogWarnFailedReadFullBytes:     "RFBTS",
	CLogWarnNoConnections:           "NOCON",
	CLogErroLoadRequestType:         "LDRQT",
	CLogErroProxyRequestType:        "PXRQT",
}

const (
	// BASE
	CLogBaseResponseModeFromService anon_logger.ILogType = iota + anon_logger.CLogFinal + 1
	CLogBaseSendNetworkMessage

	// INFO
	CLogInfoResponseFromService
	CLogInfoRecvNetworkMessage

	// WARN
	CLogWarnRequestToService
	CLogWarnUndefinedService
	CLogWarnInvalidRequestMethod
	CLogWarnFailedReadFullBytes
	CLogWarnNoConnections

	// ERRO
	CLogErroLoadRequestType
	CLogErroProxyRequestType
)
