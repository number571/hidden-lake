package http

type ILogBuilder interface {
	Build() ILogGetter
	WithMessage(string) ILogBuilder
}

type ILogGetter interface {
	GetService() string
	GetMethod() string
	GetPath() string
	GetConn() string
	GetMessage() string
}

const (
	CLogSuccess    = "_"
	CLogMethod     = "method"
	CLogDecodeBody = "decode_body"
	CLogRedirect   = "redirect"
)
