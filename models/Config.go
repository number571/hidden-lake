package models

type Config struct {
	Host Host `json:"host"`
}

type Host struct {
	Http Http `json:"http"`
	Tcp  Tcp  `json:"tcp"`
	Tls  Tls  `json:"tls"`
}

type Http struct {
	Ipv4 string `json:"ipv4"`
	Port string `json:"port"`
}

type Tls struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
}

type Tcp struct {
	Ipv4 string `json:"ipv4"`
	Port string `json:"port"`
}
