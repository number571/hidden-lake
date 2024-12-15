package anon

import (
	"fmt"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/logger"

	anon_logger "github.com/number571/go-peer/pkg/anonymity/logger"
)

const (
	cLogTemplate = "service=%s type=%s hash=%08X...%08X addr=%08X...%08X proof=%010d size=%04dB"
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
	addr := make([]byte, hashing.CHasherSize)
	if x := pLogGetter.GetPubKey(); x != nil {
		addr = hashing.NewHasher(x.ToBytes()).ToBytes()
	}

	hash := make([]byte, hashing.CHasherSize)
	if x := pLogGetter.GetHash(); x != nil {
		copy(hash, x)
	}

	return fmt.Sprintf(
		cLogTemplate,
		pLogGetter.GetService(),
		logStrType,
		hash[:4], hash[len(hash)-4:],
		addr[:4], addr[len(addr)-4:],
		pLogGetter.GetProof(),
		pLogGetter.GetSize(),
	)
}
