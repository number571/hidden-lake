// nolint:goconst
package alg

type IAlgType int

const (
	CUnknownAlg IAlgType = iota
	CQBPAlg
	CSilentAlg
)

func (p IAlgType) String() string {
	switch p {
	case CQBPAlg:
		return "qbp"
	case CSilentAlg:
		return "silent"
	default:
		return "<nil>"
	}
}

func GetAnonymityType(v string) IAlgType {
	switch v {
	case "", "qbp":
		return CQBPAlg
	case "silent":
		return CSilentAlg
	}
	return CUnknownAlg
}
