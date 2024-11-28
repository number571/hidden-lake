package http

import (
	"fmt"

	"github.com/number571/go-peer/pkg/logger"
)

const (
	cLogTemplate = "service=%s method=%s path=%s conn=%s message=%s"
)

func GetLogFunc() logger.ILogFunc {
	return func(pLogArg logger.ILogArg) string {
		logGetter, ok := pLogArg.(ILogGetter)
		if !ok {
			panic("got invalid log arg")
		}
		return fmt.Sprintf(
			cLogTemplate,
			logGetter.GetService(),
			logGetter.GetMethod(),
			logGetter.GetPath(),
			logGetter.GetConn(),
			logGetter.GetMessage(),
		)
	}
}
