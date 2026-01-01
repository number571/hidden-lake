package anon

import (
	"fmt"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/logger"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/qb/logger"
)

const (
	cLogTemplate     = "service=%s type=%s hash=%08X...%08X proof=%010d size=%04dB"
	cLogAddrTemplate = " addr=%s...%s"
	cLogConnTemplate = " conn=%s"
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
	hash := make([]byte, hashing.CHasherSize)
	if x := pLogGetter.GetHash(); x != nil {
		copy(hash, x)
	}

	log := strings.Builder{}
	log.Grow(1 << 10)

	log.WriteString(fmt.Sprintf(
		cLogTemplate,
		pLogGetter.GetService(),
		logStrType,
		hash[:4], hash[len(hash)-4:],
		pLogGetter.GetProof(),
		pLogGetter.GetSize(),
	))

	if x := pLogGetter.GetConn(); x != "" {
		log.WriteString(fmt.Sprintf(cLogConnTemplate, x))
	}
	if x := pLogGetter.GetPubKey(); x != nil {
		addr := strings.ToUpper(x.GetHasher().ToString())
		log.WriteString(fmt.Sprintf(cLogAddrTemplate, addr[:8], addr[len(addr)-8:]))
	}

	return log.String()
}
