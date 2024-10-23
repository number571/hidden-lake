package settings

type SContainer struct {
	FAliasName string `json:"alias_name"`
	FPldHead   uint64 `json:"pld_head"`
	FHexData   string `json:"hex_data"`
}
