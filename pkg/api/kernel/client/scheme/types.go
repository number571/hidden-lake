// nolint:goconst
package scheme

type ISchemeType int

const (
	CHybridScheme ISchemeType = 1 + iota
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
	return 0
}
