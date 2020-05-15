package models

type State struct {
	UsedF2F bool `json:"used_f2f"` // friend-to-friend
	UsedFSH bool `json:"used_fsh"` // file sharing
	UsedGCH bool `json:"used_gch"` // group chat
}
