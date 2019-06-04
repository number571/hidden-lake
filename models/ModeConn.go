package models

type ModeConn int8
const (
    WAIT ModeConn = -1
    NONE ModeConn =  0
    CONN ModeConn =  1
)
