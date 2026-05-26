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
