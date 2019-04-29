package models

type Connection struct {
	TempConnect string
	TempArchive []string
	NodeAddress map[string]string
	NodeLogin map[string]string
	DefaultConnections []string
}
