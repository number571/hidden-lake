package settings

import (
	"strings"

	"github.com/number571/hidden-lake/internal/adapters"
)

var (
	gAppShortNameFMT = strings.ToUpper(CAppShortName)
)

func GetAppShortNameFMT(suf ...string) string {
	return gAppShortNameFMT + strings.Join(suf, "")
}

const (
	CAppShortName = adapters.CAdapterShortName + "-" + CAppAdapterName
	CAppFullName  = adapters.CAdapterFullName + "=" + CAppAdapterName
)

const (
	CAppAdapterName = "meshtastic"
	CAppDescription = "adapts HL traffic over Meshtastic/LoRa protocol"
)

const (
	CPathYML  = CAppShortName + ".yml"
	CPathDB   = CAppShortName + ".db"
	CPathPy   = CAppShortName + ".py"
	CPathVenv = CAppShortName + ".venv"
	CPathTxt  = CAppShortName + ".txt"
)
