package models

type File struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Path string `json:"path"`
	Size uint64 `json:"size"`
}
