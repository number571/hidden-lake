package models

type From struct {
	// Login string
	Hash string
	Address string
}

type To struct {
	Hash string
	Address string
}

type Head struct {
	Title string
	Mode string
}

type PackageTCP struct {
    From From
    To To
    Head Head
    Body string
}
