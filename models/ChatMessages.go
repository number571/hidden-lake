package models

type ChatMessages struct {
	GlobalMessages []string
	LocalMessages map[string][]string
}
