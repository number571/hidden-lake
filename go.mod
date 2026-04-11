module github.com/number571/hidden-lake

go 1.24.5

require github.com/number571/go-peer v1.7.15

require (
	github.com/cloudflare/circl v1.6.3 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/number571/hidden-lake/projects/hl-client => ./projects/hl-client
