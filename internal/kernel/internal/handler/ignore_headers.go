package handler

import (
	"strings"

	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
)

type sIgnoreHeader struct {
	fHeader string
	fFilter func(string) bool
}

var gIgnoreHeaders = initIgnoreHeaders()

func initIgnoreHeaders() map[string]func(string) bool {
	headers := map[string]func(string) bool{
		hlk_settings.CHeaderResponseMode: func(_ string) bool { return true },    // ignore HLK header
		"Date":                           func(_ string) bool { return true },    // ignore due to deanonymization
		"Content-Length":                 func(_ string) bool { return true },    // ignore redundant header
		"Content-Type":                   func(v string) bool { return v == "" }, // ignore if value is null
	}

	// add headers in lower-case format
	lcHeaders := make([]sIgnoreHeader, 0, len(headers))
	for k, v := range headers {
		lcHeaders = append(
			lcHeaders,
			sIgnoreHeader{fHeader: strings.ToLower(k), fFilter: v},
		)
	}

	for _, v := range lcHeaders {
		headers[v.fHeader] = v.fFilter
	}

	return headers
}
