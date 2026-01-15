package settings

type SCommandExecRequest struct {
	FPassword string   `json:"password"`
	FCommand  []string `json:"command"`
}
