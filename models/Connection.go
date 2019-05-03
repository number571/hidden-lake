package models

type NodesF2F struct {
	NodeAddressF2F map[string]string
}

type Nodes struct {
	NodeConnection map[string]int8
	NodeAddress map[string]string
	NodeLogin map[string]string
	NodesF2F
}

type Temps struct {
	TempConnect string
	TempArchive []string
}

type Default struct {
	DefaultConnections []string
}

type Connection struct {
	Default
	Temps
	Nodes
}
