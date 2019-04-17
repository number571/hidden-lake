package models

type Messages struct {
	CurrentIdGlobal uint16
	CurrentIdLocal map[string]uint16
	NewDataExistGlobal chan bool
	NewDataExistLocal map[string]chan bool
}
