package anon

import (
	"fmt"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/logger"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
)

const (
	cLogTemplate      = "service=%s type=%s"
	cLogHashTemplate  = " hash=%08X...%08X"
	cLogProofTemplate = " proof=%d"
	cLogSizeTemplate  = " size=%dB"
	cLogConnTemplate  = " conn=%s"
)

func GetLogFunc() logger.ILogFunc {
	return func(pLogArg logger.ILogArg) string {
		logGetter, ok := pLogArg.(anon_logger.ILogGetter)
		if !ok {
			panic("got invalid log arg")
		}

		logType := logGetter.GetType()
		if logType == 0 {
			panic("got invalid log type")
		}

		logStrType, ok := gLogMap[logType]
		if !ok || logStrType == "" {
			panic("value not found in log map")
		}

		return getLog(logStrType, logGetter)
	}
}

func getLog(logStrType string, pLogGetter anon_logger.ILogGetter) string {
	log := strings.Builder{}
	log.Grow(1 << 10)

	log.WriteString(fmt.Sprintf(
		cLogTemplate,
		pLogGetter.GetService(),
		logStrType,
	))

	if x := pLogGetter.GetHash(); x != nil {
		hash := make([]byte, hashing.CHasherSize)
		copy(hash, x)
		log.WriteString(fmt.Sprintf(cLogHashTemplate, x[:4], hash[len(hash)-4:]))
	}
	if x := pLogGetter.GetProof(); x != 0 {
		log.WriteString(fmt.Sprintf(cLogProofTemplate, x))
	}
	if x := pLogGetter.GetSize(); x != 0 {
		log.WriteString(fmt.Sprintf(cLogSizeTemplate, x))
	}
	if x := pLogGetter.GetConn(); x != "" {
		log.WriteString(fmt.Sprintf(cLogConnTemplate, x))
	}

	return log.String()
}
