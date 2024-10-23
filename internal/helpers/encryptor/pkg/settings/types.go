package settings

import (
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

type SContainer struct {
	hls_settings.SPubKey
	FPldHead uint64 `json:"pld_head"`
	FHexData string `json:"hex_data"`
}
