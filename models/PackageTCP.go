package models

type From struct {
	Address string
	Login string
	Name string
}

type Head struct {
	Title string
	Mode string
}

type PackageTCP struct {
    From From
    To string
    Head Head
    Body string
}
