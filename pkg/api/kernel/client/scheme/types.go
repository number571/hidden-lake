// nolint:goconst
package scheme

type ISchemeType int

const (
	CUnknownScheme ISchemeType = iota
	CHybridScheme
	CSymmetricScheme
)

func (p ISchemeType) String() string {
	switch p {
	case CHybridScheme:
		return "hybrid"
	case CSymmetricScheme:
		return "symmetric"
	default:
		return "<nil>"
	}
}

func GetCryptoSchemeType(v string) ISchemeType {
	switch v {
	case "", "hybrid":
		return CHybridScheme
	case "symmetric":
		return CSymmetricScheme
	}
	return CUnknownScheme
}
