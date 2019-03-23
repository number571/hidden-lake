package models

type Connection struct {
	TempConnect string
	TempProfile []string
	TempArchive []string
	NodeAddress map[string]string
	Connections []string
}
